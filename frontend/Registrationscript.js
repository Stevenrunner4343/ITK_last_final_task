const API_URL = 'http://localhost:8082';

function switchTab(tab) {
    document.getElementById('signInForm').classList.toggle('hidden', tab === 'signUp');
    document.getElementById('signUpForm').classList.toggle('hidden', tab === 'signIn');
    document.getElementById('tabSignIn').classList.toggle('active', tab === 'signIn');
    document.getElementById('tabSignUp').classList.toggle('active', tab === 'signUp');
    hideMessage();
}

function showMessage(text, type) {
    const el = document.getElementById('message');
    el.textContent = text;
    el.className = 'message ' + type;
}

function hideMessage() {
    document.getElementById('message').className = 'message';
}

function showToken(token) {
    const el = document.getElementById('tokenDisplay');
    el.textContent = token;
    el.className = 'token-box show';
}

async function signIn(e) {
    e.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const res = await fetch(`${API_URL}/singIn`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        });
        const data = await res.json();
        
        if (res.ok) {
            showMessage('✅ Успешный вход!', 'success');
            if (data.token) showToken(data.token);
        } else {
            showMessage('❌ ' + (data.error || 'Ошибка входа'), 'error');
        }
    } catch (err) {
        showMessage('❌ Ошибка соединения с сервером', 'error');
    }
}

async function signUp(e) {
    e.preventDefault();
    const username = document.getElementById('regUsername').value;
    const email = document.getElementById('regEmail').value;
    const password = document.getElementById('regPassword').value;

    try {
        const res = await fetch(`${API_URL}/signUp`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, email, password })
        });
        const data = await res.json();
        
        if (res.ok) {
            showMessage('✅ Регистрация успешна!', 'success');
            if (data.token) showToken(data.token);
            setTimeout(() => switchTab('signIn'), 1500);
        } else {
            showMessage('❌ ' + (data.error || 'Ошибка регистрации'), 'error');
        }
    } catch (err) {
        showMessage('❌ Ошибка соединения с сервером', 'error');
    }
}