FROM python:3.14-alpine
WORKDIR /pytest
COPY pytest/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY ./pytest/*.py .
RUN ls ./
CMD ["pytest", "-v", "-s"]
