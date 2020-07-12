# Prepare deployment package

## Prerequisite

1.You have python virtual environment in your instance.

## Pack your dependency

```bash
chmod +x pack.sh && ./pack.sh
```

## Want to support more package

```bash
source venv/bin/activate
pip install <new_package>
pip freeze > requirements.txt
```
