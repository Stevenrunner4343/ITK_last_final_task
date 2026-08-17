console.log('Скрипт загружен!');

async function sendData() {
    const data = {
        Id: parseInt(document.getElementById('userId').value),
        Operation: document.getElementById('operation').value,
        Amount: parseInt(document.getElementById('amount').value)
    };

    const response = await fetch('http://localhost:8082/operation', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(data)
    });
    console.log('POST отправлен!', data);
}

async function getBalance() {
    const userId = parseInt(document.getElementById('balanceUserId').value);
    
    const data = {
        Id: userId
    };

    const response = await fetch('http://localhost:8082/balance', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(data)
    });
    
    const result = await response.json();
    console.log('GET получен!', result);
    
    if (result.error) {
        document.getElementById('balanceResult').textContent = '❌ ' + result.error;

    }else{
        document.getElementById('balanceResult').textContent = '💰 Balance: ' + result;
    }
            
}


document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('sendBtn').addEventListener('click', sendData);
    document.getElementById('balanceBtn').addEventListener('click', getBalance);
})
;