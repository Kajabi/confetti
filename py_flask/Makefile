.PHONY : setup venv run

run: venv setup
	./venv/bin/python -m flask run --port 3000

venv:
	python3 -m venv venv

setup:
	./venv/bin/python -m pip install -r requirements.txt
