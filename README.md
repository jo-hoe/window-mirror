# Window Mirror

Mirrors Windows on Windows OS horizontally.

## Prerequisites to run locally

## Python

This program is tested on python 3.13.

## Run with python

Install the requirements by.

```PowerShell
pip install -r requirements.txt
```

and run the program by

```PowerShell
python main.py
```

## Building an .exe

You can also create an `.exe` by installing

```PowerShell
pip install installer
```

and running

```PowerShell
pyinstaller -F main.py
```

The `.exe` can then be found in the newly created `dist` folder.

## Run with "make"

The project is using `make`.
`make` is not strictly required, but it helps and documents commonly used commands.

You can directly install it from [gnuwin32](https://gnuwin32.sourceforge.net/packages/make.htm) or via `winget`

```PowerShell
winget install GnuWin32.Make
```

You will also need Docker and Python.
Run `make init` to install all dependencies in a virtual Python environment.

### How to Use

You can check all `make` commands by running.

```bash
make help
```

## How To Run Locally

For local development and testing, put a [`config.yaml`](#configuration) in folder `dev`.
To run the service locally, you can use `docker-compose` or just run it via make:

```bash
make start
```
