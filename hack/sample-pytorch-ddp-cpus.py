import torch
import torch.nn.functional as F
from torch.utils.data import DistributedSampler
from torchvision import datasets, transforms

# [1] Setup PyTorch DDP. WORLD_SIZE and RANK environments will be set by Training Operator.
torch.distributed.init_process_group(backend="gloo")
Distributor = torch.nn.parallel.DistributedDataParallel


# [2] Create PyTorch CNN Model.
class Net(torch.nn.Module):

    def __init__(self):
        super(Net, self).__init__()
        self.conv1 = torch.nn.Conv2d(1, 20, 5, 1)
        self.conv2 = torch.nn.Conv2d(20, 50, 5, 1)
        self.fc1 = torch.nn.Linear(4 * 4 * 50, 500)
        self.fc2 = torch.nn.Linear(500, 10)

    def forward(self, x):
        x = F.relu(self.conv1(x))
        x = F.max_pool2d(x, 2, 2)
        x = F.relu(self.conv2(x))
        x = F.max_pool2d(x, 2, 2)
        x = x.view(-1, 4 * 4 * 50)
        x = F.relu(self.fc1(x))
        x = self.fc2(x)
        return F.log_softmax(x, dim=1)


# [3] Attach model to CPU and distributor.
device = "cpu"
model = Net().to(device)
model = Distributor(model)
optimizer = torch.optim.SGD(model.parameters(), lr=0.01, momentum=0.5)

# [4] Setup FashionMNIST dataloader and distribute data across PyTorchJob workers.
dataset = datasets.FashionMNIST(
    "./data",
    download=True,
    train=True,
    transform=transforms.Compose([transforms.ToTensor()]),
)
train_loader = torch.utils.data.DataLoader(
    dataset=dataset,
    batch_size=128,
    sampler=DistributedSampler(dataset),
)

# [5] Start model Training.
for epoch in range(3):
    for batch_idx, (data, target) in enumerate(train_loader):
        # Attach Tensors to the device.
        data = data.to(device)
        target = target.to(device)

        optimizer.zero_grad()
        output = model(data)
        loss = F.nll_loss(output, target)
        loss.backward()
        optimizer.step()
        if batch_idx % 10 == 0:
            print("Train Epoch: {} [{}/{} ({:.0f}%)]\tloss={:.4f}".format(
                epoch,
                batch_idx * len(data),
                len(train_loader.dataset),
                100.0 * batch_idx / len(train_loader),
                loss.item(),
            ))
