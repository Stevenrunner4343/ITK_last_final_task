document.getElementById('registerBtn').addEventListener('click', async function() {
    const data = {
        Username: document.getElementById('username').value,
        Email: document.getElementById('email').value,
        Password: document.getElementById('password').value
    };
    
    try {
        const response = await fetch('http://localhost:8082/signUp', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(data)
        });
        
        const result = await response.json();
        
        if (response.status === 201) {
            document.getElementById('message').className = 'green';
            document.getElementById('message').textContent = '✅ ' + result.message;
            
            // Переход на wallet через 1 секунду
            setTimeout(() => {
                window.location.href = 'wallet.html';
            }, 1000);
        } else {
            document.getElementById('message').className = 'red';
            document.getElementById('message').textContent = '❌ ' + result.error;
        }
    } catch (error) {
        document.getElementById('message').className = 'red';
        document.getElementById('message').textContent = '❌ Ошибка сервера';
    }
});