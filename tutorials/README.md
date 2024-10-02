# Tutorial

This directory contains Jupyter Notebooks for the LLMariner Tutorial.

## Setting up Juypter

Please follow https://jupyter.org/install to install Juypter. The following is an example
install procedure

```bash
VENV_PATH=<your virtual env path>
python -m venv "${VENV_PATH}"
source "${VENV_PATH}/bin/activate"
pip install jupyterlab
```

Once installed, you can launch JupyterLab with:

```bash
jupyter lab
```

## Internal Notes

### Converting `.md` Files to `.ipynb` Files

We use [jupytext](https://jupytext.readthedocs.io/en/latest/) to do the conversion between
Markdown files Jupyter Notebooks. Run the following command to install `jupytext`:

```bash
pip install jupytext
jupyter labextension install jupyterlab-jupytext@1.2.2
```

If file syncing is not happening, you can run `jupytext --sync <markdown file name>` to force it.
