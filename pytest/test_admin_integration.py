"""
Интеграционные тесты для сервиса Anti-Bruteforce с полной проверкой административных функций.

Тесты проверяют:
1. Авторизацию администратора (корректный и некорректный API ключ)
2. Rate limiting для логина, пароля и IP
3. Управление blacklist (добавление, получение, удаление)
4. Управление whitelist (добавление, получение, удаление)
5. Сброс bucket'ов
6. Приоритет whitelist и blacklist над rate limiting
7. Проверку доступа к административным эндпоинтам без авторизации
"""

import requests
import time
from http import HTTPStatus
import pytest


BASE_URL = 'http://10.0.0.3:8080'
# API ключ должен совпадать с конфигурацией в conf.yaml или переменной окружения
ADMIN_API_KEY = 'test-admin-key-123'  # Замените на реальный ключ из конфигурации


class TestAdminAuthentication:
    """Тесты аутентификации администратора"""
    
    def test_admin_access_with_valid_key(self):
        """Проверка доступа к административному API с корректным ключом"""
        url = f'{BASE_URL}/admin/auth/whitelist'
        headers = {
            'Content-Type': 'application/json',
            'X-Admin-Key': ADMIN_API_KEY
        }
        
        response = requests.get(url, headers=headers)
        print(f'Admin access with valid key: status={response.status_code}')
        
        # Должен быть успешный доступ (200 OK)
        assert response.status_code == HTTPStatus.OK, \
            f"Expected 200, got {response.status_code}"
    
    def test_admin_access_without_key(self):
        """Проверка отказа в доступе без API ключа"""
        url = f'{BASE_URL}/admin//auth/whitelist'
        headers = {'Content-Type': 'application/json'}
        
        response = requests.get(url, headers=headers)
        print(f'Admin access without key: status={response.status_code}')
        
        # Должен быть отказ в доступе (401 Unauthorized)
        assert response.status_code == HTTPStatus.UNAUTHORIZED, \
            f"Expected 401, got {response.status_code}"
    
    def test_admin_access_with_invalid_key(self):
        """Проверка отказа в доступе с неверным API ключом"""
        url = f'{BASE_URL}/admin/auth/whitelist'
        headers = {
            'Content-Type': 'application/json',
            'X-Admin-Key': 'invalid-key-wrong'
        }
        
        response = requests.get(url, headers=headers)
        print(f'Admin access with invalid key: status={response.status_code}')
        
        # Должен быть отказ в доступе (401 Unauthorized)
        assert response.status_code == HTTPStatus.UNAUTHORIZED, \
            f"Expected 401, got {response.status_code}"
    
    def test_all_admin_endpoints_require_auth(self):
        """Проверка, что все административные эндпоинты требуют аутентификации"""
        endpoints = [
            ('GET', '/admin/auth/whitelist'),
            ('POST', '/admin//auth/whitelist'),
            ('DELETE', '/admin/auth/whitelist'),
            ('GET', '/admin/auth/blacklist'),
            ('POST', '/admin/auth/blacklist'),
            ('DELETE', '/admin/auth/blacklist'),
            ('DELETE', '/admin/auth/reset'),
        ]
        
        headers = {'Content-Type': 'application/json'}
        
        for method, endpoint in endpoints:
            url = f'{BASE_URL}{endpoint}'
            
            if method == 'GET':
                response = requests.get(url, headers=headers)
            elif method == 'POST':
                response = requests.post(url, json={}, headers=headers)
            elif method == 'DELETE':
                response = requests.delete(url, json={}, headers=headers)
            
            print(f'{method} {endpoint}: status={response.status_code}')
            assert response.status_code == HTTPStatus.UNAUTHORIZED, \
                f"{method} {endpoint} should require authentication"


class TestBlacklistManagement:
    """Тесты управления черным списком"""
    
    def setup_method(self):
        """Очистка blacklist перед каждым тестом"""
        self.admin_headers = {
            'Content-Type': 'application/json',
            'X-Admin-Key': ADMIN_API_KEY
        }
        self._clear_blacklist()
    
    def teardown_method(self):
        """Очистка blacklist после каждого теста"""
        self._clear_blacklist()
    
    def _clear_blacklist(self):
        """Вспомогательный метод для очистки blacklist"""
        url = f'{BASE_URL}/admin/auth/blacklist'
        response = requests.get(url, headers=self.admin_headers)
        if response.status_code == HTTPStatus.OK:
            for ip_network in response.json():
                requests.delete(url, json=ip_network, headers=self.admin_headers)
    
    def test_add_ip_to_blacklist(self):
        """Тест добавления IP в blacklist"""
        url = f'{BASE_URL}/admin/auth/blacklist'
        data = {
            'ip': '192.168.1.0',
            'mask': '255.255.255.0'
        }
        
        response = requests.post(url, json=data, headers=self.admin_headers)
        print(f'Add to blacklist: status={response.status_code}')
        
        assert response.status_code == HTTPStatus.NO_CONTENT, \
            f"Expected 204, got {response.status_code}"
    
    def test_add_duplicate_ip_to_blacklist(self):
        """Тест добавления дубликата IP в blacklist"""
        url = f'{BASE_URL}/admin/auth/blacklist'
        data = {
            'ip': '192.168.2.0',
            'mask': '255.255.255.0'
        }
        
        # Первое добавление
        response1 = requests.post(url, json=data, headers=self.admin_headers)
        assert response1.status_code == HTTPStatus.NO_CONTENT
        
        # Повторное добавление
        response2 = requests.post(url, json=data, headers=self.admin_headers)
        print(f'Add duplicate to blacklist: status={response2.status_code}')
        
        assert response2.status_code == HTTPStatus.BAD_REQUEST, \
            "Duplicate IP should return 400"
    
    def test_get_blacklist(self):
        """Тест получения списка IP из blacklist"""
        url = f'{BASE_URL}/admin/auth/blacklist'
        
        # Добавляем несколько IP
        test_networks = [
            {'ip': '10.0.0.0', 'mask': '255.0.0.0'},
            {'ip': '172.16.0.0', 'mask': '255.240.0.0'},
        ]
        
        for network in test_networks:
            requests.post(url, json=network, headers=self.admin_headers)
        
        # Получаем список
        response = requests.get(url, headers=self.admin_headers)
        print(f'Get blacklist: status={response.status_code}, items={len(response.json())}')
        
        assert response.status_code == HTTPStatus.OK
        assert len(response.json()) >= 2, "Should contain at least 2 networks"
    
    def test_delete_ip_from_blacklist(self):
        """Тест удаления IP из blacklist"""
        url = f'{BASE_URL}/admin/auth/blacklist'
        data = {
            'ip': '192.168.3.0',
            'mask': '255.255.255.0'
        }
        
        # Добавляем
        requests.post(url, json=data, headers=self.admin_headers)
        
        # Удаляем
        response = requests.delete(url, json=data, headers=self.admin_headers)
        print(f'Delete from blacklist: status={response.status_code}')
        
        assert response.status_code == HTTPStatus.NO_CONTENT, \
            f"Expected 204, got {response.status_code}"
    
    def test_blacklist_blocks_requests(self):
        """Тест блокировки запросов из blacklist"""
        blacklist_url = f'{BASE_URL}/admin/auth/blacklist'
        auth_url = f'{BASE_URL}/auth/check'
        
        # Добавляем IP в blacklist
        blacklist_data = {
            'ip': '12.0.0.0',
            'mask': '255.255.0.0'
        }
        requests.post(blacklist_url, json=blacklist_data, headers=self.admin_headers)
        
        # Пытаемся авторизоваться с IP из blacklist
        auth_data = {
            'login': 'user1',
            'password': 'pass123',
            'ip': '12.0.1.1'  # Попадает в подсеть 12.0.0.0/16
        }
        
        response = requests.post(auth_url, json=auth_data)
        print(f'Auth from blacklisted IP: response={response.content}')
        
        assert response.content == b'ok=false', \
            "Request from blacklisted IP should be blocked"


class TestWhitelistManagement:
    """Тесты управления белым списком"""
    
    def setup_method(self):
        """Очистка whitelist перед каждым тестом"""
        self.admin_headers = {
            'Content-Type': 'application/json',
            'X-Admin-Key': ADMIN_API_KEY
        }
        self._clear_whitelist()
    
    def teardown_method(self):
        """Очистка whitelist после каждого теста"""
        self._clear_whitelist()
    
    def _clear_whitelist(self):
        """Вспомогательный метод для очистки whitelist"""
        url = f'{BASE_URL}/admin/auth/whitelist'
        response = requests.get(url, headers=self.admin_headers)
        if response.status_code == HTTPStatus.OK:
            for ip_network in response.json():
                requests.delete(url, json=ip_network, headers=self.admin_headers)
    
    def test_add_ip_to_whitelist(self):
        """Тест добавления IP в whitelist"""
        url = f'{BASE_URL}/admin/auth/whitelist'
        data = {
            'ip': '10.10.0.0',
            'mask': '255.255.0.0'
        }
        
        response = requests.post(url, json=data, headers=self.admin_headers)
        print(f'Add to whitelist: status={response.status_code}')
        
        assert response.status_code == HTTPStatus.NO_CONTENT, \
            f"Expected 204, got {response.status_code}"
    
    def test_get_whitelist(self):
        """Тест получения списка IP из whitelist"""
        url = f'{BASE_URL}/admin/auth/whitelist'
        
        # Добавляем несколько IP
        test_networks = [
            {'ip': '192.168.0.0', 'mask': '255.255.0.0'},
            {'ip': '172.20.0.0', 'mask': '255.255.240.0'},
        ]
        
        for network in test_networks:
            requests.post(url, json=network, headers=self.admin_headers)
        
        # Получаем список
        response = requests.get(url, headers=self.admin_headers)
        print(f'Get whitelist: status={response.status_code}, items={len(response.json())}')
        
        assert response.status_code == HTTPStatus.OK
        assert len(response.json()) >= 2, "Should contain at least 2 networks"
    
    def test_delete_ip_from_whitelist(self):
        """Тест удаления IP из whitelist"""
        url = f'{BASE_URL}/admin/auth/whitelist'
        data = {
            'ip': '192.168.10.0',
            'mask': '255.255.255.0'
        }
        
        # Добавляем
        requests.post(url, json=data, headers=self.admin_headers)
        
        # Удаляем
        response = requests.delete(url, json=data, headers=self.admin_headers)
        print(f'Delete from whitelist: status={response.status_code}')
        
        assert response.status_code == HTTPStatus.NO_CONTENT, \
            f"Expected 204, got {response.status_code}"
    
    def test_whitelist_allows_unlimited_requests(self):
        """Тест разрешения неограниченных запросов из whitelist"""
        whitelist_url = f'{BASE_URL}/admin/auth/whitelist'
        auth_url = f'{BASE_URL}/auth/check'
        
        # Добавляем IP в whitelist
        whitelist_data = {
            'ip': '20.0.0.0',
            'mask': '255.255.0.0'
        }
        requests.post(whitelist_url, json=whitelist_data, headers=self.admin_headers)
        
        # Делаем много запросов (больше лимита)
        auth_data = {
            'login': 'spammer',
            'password': 'pass123',
            'ip': '20.0.1.1'  # Попадает в подсеть 20.0.0.0/16
        }
        
        # Делаем 20 запросов (больше чем лимит логина = 10)
        for i in range(20):
            response = requests.post(auth_url, json=auth_data)
            print(f'Request {i+1} from whitelisted IP: response={response.content}')
            assert response.content == b'ok=true', \
                f"Request {i+1} from whitelist should always be allowed"


class TestBucketReset:
    """Тесты сброса bucket'ов"""
    
    def setup_method(self):
        """Подготовка перед каждым тестом"""
        self.admin_headers = {
            'Content-Type': 'application/json',
            'X-Admin-Key': ADMIN_API_KEY
        }
    
    def test_reset_bucket_for_login_and_ip(self):
        """Тест сброса bucket для логина и IP"""
        auth_url = f'{BASE_URL}/auth/check'
        reset_url = f'{BASE_URL}/admin/auth/reset'
        
        test_login = 'reset_test_user'
        test_ip = '50.0.0.1'
        
        # Исчерпываем лимит для логина (10 попыток)
        for i in range(10):
            data = {
                'login': test_login,
                'password': f'pass{i}',
                'ip': test_ip
            }
            response = requests.post(auth_url, json=data)
            print(f'Attempt {i+1}: {response.content}')
        
        # 11-я попытка должна быть заблокирована
        data = {
            'login': test_login,
            'password': 'pass11',
            'ip': test_ip
        }
        response = requests.post(auth_url, json=data)
        assert response.content == b'ok=false', "Should be blocked after limit"
        
        # Сбрасываем bucket
        reset_data = {
            'login': test_login,
            'ip': test_ip
        }
        response = requests.delete(reset_url, json=reset_data, headers=self.admin_headers)
        print(f'Reset bucket: status={response.status_code}, content={response.content}')
        
        assert response.status_code == HTTPStatus.OK
        assert b'resetLogin=true' in response.content
        assert b'resetIp=true' in response.content
        
        # После сброса должна пройти новая попытка
        response = requests.post(auth_url, json=data)
        print(f'After reset: {response.content}')
        assert response.content == b'ok=true', "Should be allowed after reset"
    
    def test_reset_bucket_requires_admin_auth(self):
        """Тест требования аутентификации для сброса bucket"""
        reset_url = f'{BASE_URL}/admin/auth/reset'
        data = {
            'login': 'someuser',
            'ip': '1.2.3.4'
        }
        
        headers = {'Content-Type': 'application/json'}
        response = requests.delete(reset_url, json=data, headers=headers)
        
        assert response.status_code == HTTPStatus.UNAUTHORIZED, \
            "Reset should require admin authentication"


class TestRateLimiting:
    """Тесты ограничения частоты запросов (Rate Limiting)"""
    
    def test_login_rate_limit(self):
        """Тест лимита по логину (10 попыток в минуту)"""
        url = f'{BASE_URL}/auth/check'
        test_login = 'rate_limit_user'
        
        # Делаем 10 попыток - все должны пройти
        for i in range(10):
            data = {
                'login': test_login,
                'password': f'pass{i}',
                'ip': f'60.0.0.{i}'  # Разные IP
            }
            response = requests.post(url, json=data)
            print(f'Login attempt {i+1}: {response.content}')
            assert response.content == b'ok=true', \
                f"Attempt {i+1} should be allowed"
        
        # 11-я попытка должна быть заблокирована
        data = {
            'login': test_login,
            'password': 'pass11',
            'ip': '60.0.0.11'
        }
        response = requests.post(url, json=data)
        print(f'Login attempt 11: {response.content}')
        assert response.content == b'ok=false', \
            "11th attempt should be blocked (login limit)"
    
    @pytest.mark.skip(reason="Too many requests cause file descriptor issues")
    def test_password_rate_limit(self):
        """Тест лимита по паролю (100 попыток)"""
        from concurrent.futures import ThreadPoolExecutor
        
        url = f'{BASE_URL}/auth/check'
        test_password = f'password_{int(time.time())}'
        
        def make_request(i):
            return requests.post(url, json={
                'login': f'user{i}_{int(time.time())}',
                'password': test_password,
                'ip': f'70.{i // 256}.{i % 256}.1'
            })
        
        # Делаем 150 запросов (на 50 больше лимита)
        start = time.time()
        with ThreadPoolExecutor(max_workers=10) as executor:
            responses = list(executor.map(make_request, range(150)))
        elapsed = time.time() - start
        
        success = sum(1 for r in responses if r.content == b'ok=true')
        blocked = sum(1 for r in responses if r.content == b'ok=false')
        
        print(f"\nPassword Rate Limit Test:")
        print(f"  Total requests: 150")
        print(f"  Time: {elapsed:.2f}s")
        print(f"  Success: {success}")
        print(f"  Blocked: {blocked}")
        print(f"  Configured limit: 100")
        
        # Должно быть много заблокированных
        assert blocked >= 40, \
            f"Should block at least 40 requests, blocked: {blocked}"
        
        # Должно пройти примерно 100 (±10 на bucket leak)
        assert 95 <= success <= 110, \
            f"Expected ~100 to pass (±10), got {success}"
        
        print(f"  ✓ Password rate limiting works correctly")
    
    @pytest.mark.skip(reason="Too many requests cause file descriptor issues")
    def test_ip_rate_limit(self):
        """Тест лимита по IP"""
        from concurrent.futures import ThreadPoolExecutor
        
        url = f'{BASE_URL}/auth/check'
        test_ip = f'80.{int(time.time()) % 255}.0.1'
        
        def make_request(i):
            return requests.post(url, json={
                'login': f'ipuser{i}_{int(time.time())}',
                'password': f'pass{i}_{int(time.time())}',
                'ip': test_ip
            })
        
        # Делаем 1500 запросов
        start = time.time()
        with ThreadPoolExecutor(max_workers=10) as executor:
            responses = list(executor.map(make_request, range(1500)))
        elapsed = time.time() - start
        
        success = sum(1 for r in responses if r.content == b'ok=true')
        blocked = sum(1 for r in responses if r.content == b'ok=false')
        
        print(f"\nIP Rate Limit Test:")
        print(f"  Total requests: 1500")
        print(f"  Time: {elapsed:.2f}s")
        print(f"  Success: {success}")
        print(f"  Blocked: {blocked}")
        print(f"  Configured limit: 1000")
        
        # Должно быть много заблокированных
        assert blocked >= 400, \
            f"Should block at least 400 requests, blocked: {blocked}"
        
        # Должно пройти примерно 1000 (±50 на bucket leak и race conditions)
        assert 980 <= success <= 1050, \
            f"Expected ~1000 to pass (±50), got {success}"
        
        print(f"  ✓ IP rate limiting works correctly")

class TestInputValidation:
    """Тесты валидации входных данных"""
    
    def test_invalid_ip_format_in_blacklist(self):
        """Тест добавления некорректного IP в blacklist"""
        url = f'{BASE_URL}/admin/auth/blacklist'
        headers = {
            'Content-Type': 'application/json',
            'X-Admin-Key': ADMIN_API_KEY
        }
        
        invalid_data = [
            {'ip': '999.999.999.999', 'mask': '255.255.255.0'},
            {'ip': 'invalid', 'mask': '255.255.255.0'},
            {'ip': '192.168.1.1', 'mask': 'invalid'},
        ]
        
        for data in invalid_data:
            response = requests.post(url, json=data, headers=headers)
            print(f'Invalid IP data: {data}, status={response.status_code}')
            assert response.status_code == HTTPStatus.BAD_REQUEST, \
                f"Invalid IP format should return 400: {data}"
    
    def test_invalid_auth_request(self):
        """Тест некорректного запроса авторизации"""
        url = f'{BASE_URL}/auth/check'
        
        invalid_requests = [
            {},  # Пустой запрос
            {'login': 'user'},  # Нет пароля и IP
            {'password': 'pass'},  # Нет логина и IP
            {'ip': '1.2.3.4'},  # Нет логина и пароля
        ]
        
        for data in invalid_requests:
            response = requests.post(url, json=data)
            print(f'Invalid auth request: {data}, status={response.status_code}')
            assert response.status_code == HTTPStatus.BAD_REQUEST, \
                f"Invalid request should return 400: {data}"


class TestPriorityAndEdgeCases:
    """Тесты приоритета и граничных случаев"""
    
    def setup_method(self):
        """Подготовка перед каждым тестом"""
        self.admin_headers = {
            'Content-Type': 'application/json',
            'X-Admin-Key': ADMIN_API_KEY
        }
    
    def test_whitelist_overrides_rate_limit(self):
        """Тест приоритета whitelist над rate limiting"""
        whitelist_url = f'{BASE_URL}/admin/auth/whitelist'
        auth_url = f'{BASE_URL}/auth/check'
        
        # Добавляем IP в whitelist
        whitelist_data = {'ip': '100.0.0.0', 'mask': '255.255.0.0'}
        requests.post(whitelist_url, json=whitelist_data, headers=self.admin_headers)
        
        test_login = 'whitelist_user'
        test_ip = '100.0.1.1'
        
        # Делаем больше запросов, чем позволяет лимит логина
        for i in range(20):
            data = {
                'login': test_login,
                'password': f'pass{i}',
                'ip': test_ip
            }
            response = requests.post(auth_url, json=data)
            assert response.content == b'ok=true', \
                f"Whitelisted IP should bypass rate limit (attempt {i+1})"
        
        # Cleanup
        requests.delete(whitelist_url, json=whitelist_data, headers=self.admin_headers)
    
    def test_blacklist_overrides_whitelist(self):
        """Тест приоритета blacklist над whitelist (если IP в обоих списках)"""
        whitelist_url = f'{BASE_URL}/admin/auth/whitelist'
        blacklist_url = f'{BASE_URL}/admin/auth/blacklist'
        auth_url = f'{BASE_URL}/auth/check'
        
        network_data = {'ip': '110.0.0.0', 'mask': '255.255.0.0'}
        
        # Добавляем в whitelist
        requests.post(whitelist_url, json=network_data, headers=self.admin_headers)
        
        # Добавляем в blacklist
        requests.post(blacklist_url, json=network_data, headers=self.admin_headers)
        
        # Пытаемся авторизоваться
        data = {
            'login': 'user',
            'password': 'pass',
            'ip': '110.0.1.1'
        }
        response = requests.post(auth_url, json=data)
        
        # Blacklist должен иметь приоритет
        assert response.content == b'ok=false', \
            "Blacklist should override whitelist"
        
        # Cleanup
        requests.delete(whitelist_url, json=network_data, headers=self.admin_headers)
        requests.delete(blacklist_url, json=network_data, headers=self.admin_headers)
    
    def test_subnet_matching(self):
        """Тест корректного определения принадлежности IP к подсети"""
        blacklist_url = f'{BASE_URL}/admin/auth/blacklist'
        auth_url = f'{BASE_URL}/auth/check'
        
        # Добавляем подсеть /24
        network_data = {'ip': '120.0.1.0', 'mask': '255.255.255.0'}
        requests.post(blacklist_url, json=network_data, headers=self.admin_headers)
        
        # IP внутри подсети должны блокироваться
        blocked_ips = ['120.0.1.1', '120.0.1.100', '120.0.1.254']
        for ip in blocked_ips:
            data = {'login': 'user', 'password': 'pass', 'ip': ip}
            response = requests.post(auth_url, json=data)
            assert response.content == b'ok=false', \
                f"IP {ip} should be blocked (inside subnet)"
        
        # IP вне подсети не должны блокироваться
        allowed_ips = ['120.0.0.1', '120.0.2.1', '121.0.1.1']
        for ip in allowed_ips:
            data = {'login': 'user', 'password': 'pass', 'ip': ip}
            response = requests.post(auth_url, json=data)
            assert response.content == b'ok=true', \
                f"IP {ip} should be allowed (outside subnet)"
        
        # Cleanup
        requests.delete(blacklist_url, json=network_data, headers=self.admin_headers)


if __name__ == '__main__':
    # Запуск тестов с подробным выводом
    pytest.main([__file__, '-v', '-s'])
