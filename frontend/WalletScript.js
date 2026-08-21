
console.log('Wallet loaded!');

// Send operation
document.getElementById('sendBtn').addEventListener('click', async function() {
    const data = {
        Id: parseInt(document.getElementById('userId').value),
        Operation: document.getElementById('operation').value,
        Amount: parseInt(document.getElementById('amount').value)
    };
    
    await fetch('http://localhost:8082/operation', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(data)
    });
    alert('✅ Operation sent!');
});

// Get balance
document.getElementById('balanceBtn').addEventListener('click', async function() {
    const userId = parseInt(document.getElementById('balanceUserId').value);
    
    const response = await fetch('http://localhost:8082/balance', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({ Id: userId })
    });
    
    const result = await response.json();
    
    if (result.error) {
        document.getElementById('balanceResult').textContent = '❌ ' + result.error;
    } else {
        document.getElementById('balanceResult').textContent = '💰 Balance: ' + result;
    }
});

// Logout
document.getElementById('logoutBtn').addEventListener('click', function() {
    window.location.href = 'registration.html';
});
