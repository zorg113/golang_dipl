import requests
import time
from http import HTTPStatus


def run_auth(login:str)-> bytes:
    url='http://10.0.0.3:8080/auth/check'
    headers = {'Content-Type': 'application/json'}
    data = {
        'login': login,
        'password': 'mypasswrd',
        'ip': '10.0.0.88'}
    
    response = requests.post(url, json=data, headers=headers)
    print(f'response status code %s  value %s'%(response.status_code,response.content)) 
    return response.content
    
def test_bucket_erase():
    print("start test server")
    for i in range(0,10):
        assert run_auth('spammer')==b'ok=true'
    assert run_auth('spammer')==b'ok=false'
    assert run_auth('spammer')==b'ok=false'
    time.sleep(60)        
    assert run_auth('spammer')==b'ok=true'

def run_auth_ip(login:str, ip:str) -> bytes:
    url='http://10.0.0.3:8080/auth/check'
    headers = {'Content-Type': 'application/json'}
    data = {
        'login': login,
        'password': 'mypasswrd',
        'ip': ip}
    
    response = requests.post(url, json=data, headers=headers)
    print(f'response status code %s  value %s'%(response.status_code,response.content)) 
    return response.content

def test_blacklist():
    url='http://10.0.0.3:8080/auth/blacklist'
    headers = {'Content-Type': 'application/json'}
    response = requests.get(url, None, headers=headers)
    for i in response.json():
        response = requests.delete(url, json=i, headers=headers)   
   
    data = [{
        'ip': '12.0.0.88',
        'mask': '255.255.0.0'
        },
        {
        'ip': '16.7.8.88',
        'mask': '255.255.255.0'
        }
    ]    
    response = requests.post(url, json=data[0], headers=headers)
    print(f'response status code %s  value %s'%(response.status_code,response.content)) 
    response = requests.post(url, json=data[1], headers=headers)
    print(f'response status code %s  value %s'%(response.status_code,response.content))
    assert run_auth_ip('spamer','12.0.1.1')==b'ok=false'
    assert run_auth_ip('spamer','16.7.9.1')==b'ok=true'
    response = requests.get(url, None, headers=headers)
    to_delete = [{
        'ip': '12.0.0.0',
        'mask': '255.255.0.0'
        },
        {
        'ip': '16.7.8.0',
        'mask': '255.255.255.0'
        }
    ] 
    for i in response.json():
        if i in to_delete:
            continue
        else:
            assert i in to_delete
    response = requests.delete(url, json=to_delete[0], headers=headers)   
    assert response.status_code==HTTPStatus.NO_CONTENT
    response = requests.delete(url, json=to_delete[1], headers=headers)   
    assert response.status_code==HTTPStatus.NO_CONTENT


def test_whitelist():
    url='http://10.0.0.3:8080/auth/whitelist'
    headers = {'Content-Type': 'application/json'}
    response = requests.get(url, None, headers=headers)
    for i in response.json():
        response = requests.delete(url, json=i, headers=headers)         
    data = [{
        'ip': '12.0.0.88',
        'mask': '255.255.0.0'
        },
        {
        'ip': '16.7.8.88',
        'mask': '255.255.255.0'
        }
    ]    
    response = requests.post(url, json=data[0], headers=headers)
    print(f'response status code %s  value %s'%(response.status_code,response.content)) 
    response = requests.post(url, json=data[1], headers=headers)
    print(f'response status code %s  value %s'%(response.status_code,response.content))
    assert run_auth_ip('spamer','12.0.1.1')==b'ok=true'
    assert run_auth_ip('spamer','16.7.9.1')==b'ok=true'
    response = requests.get(url, None, headers=headers)
    to_delete = [{
        'ip': '12.0.0.0',
        'mask': '255.255.0.0'
        },
        {
        'ip': '16.7.8.0',
        'mask': '255.255.255.0'
        }
    ] 
    for i in response.json():
        if i in to_delete:
            continue
        else:
            assert i in to_delete
    response = requests.delete(url, json=to_delete[0], headers=headers)   
    assert response.status_code==HTTPStatus.NO_CONTENT
    response = requests.delete(url, json=to_delete[1], headers=headers)   
    assert response.status_code==HTTPStatus.NO_CONTENT


def test_witelist():
    pass



# if __name__ == '__main__':
#     main()