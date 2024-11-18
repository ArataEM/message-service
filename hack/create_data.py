from faker import Faker
import random
import requests
from uuid_extensions import uuid7str

def main():
    fake = Faker()

    users = []
    for i in range(100):
        users.append(uuid7str())

    for i in range(10000):
        data = {
            "user_id": random.choice(users),
            "text": fake.text()
        }
        requests.post("http://localhost:8080/messages", json=data)


if __name__ == '__main__':
    main()
