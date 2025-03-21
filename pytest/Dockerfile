FROM python:3.10-alpine
WORKDIR /pytest
RUN pip install requests pytest
COPY ./pytest/test_blocked.py ./
RUN ls ./
CMD ["pytest", "test_blocked.py"]
