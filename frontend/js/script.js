let token = localStorage.getItem('token');
const authContainer = document.getElementById('auth-container');
const mainContainer = document.getElementById('main-container');
if (token){
    authContainer.style.display = "none";
    mainContainer.style.display = "block";
} else {
    authContainer.style.display = "block";
    mainContainer.style.display = "none";
}
document.getElementById('signupBtn').addEventListener('click', () => {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    fetch('http://172.19.0.3:8080/signup', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
    })
    .then(response => response.json())
    .then(data => {
        localStorage.setItem('token', data.token);
        document.getElementById('auth-message').textContent = "サインアップ完了";
        authContainer.style.display = "none";
        mainContainer.style.display = "block";
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('auth-message').textContent = "サインアップに失敗";
    });
});

document.getElementById('loginBtn').addEventListener('click', () => {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    fetch('http://172.19.0.3:8080/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
    })
    .then(response => response.json())
    .then(data => {
        localStorage.setItem('token', data.token);
        document.getElementById('auth-message').textContent = "ログイン完了";
        authContainer.style.display = "none";
        mainContainer.style.display = "block";
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('auth-message').textContent = "ログインに失敗";
    });
});


document.getElementById('getLocationBtn').addEventListener('click', () => {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(success, error);
    } else {
        document.getElementById('status').textContent = "このブラウザは位置情報をサポートしていません。";
    }
});

document.getElementById('putLocationBtn').addEventListener('click', () => {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(putUserLocation, error);
    } else {
        document.getElementById('status').textContent = "このブラウザは位置情報をサポートしていません。";
    }
})

function success(position) {
    const latitude = position.coords.latitude;
    const longitude = position.coords.longitude;
    const human_sensor = true; 
    const light_sensor = true;
    document.getElementById('status').textContent = "位置情報を送信中...";

    fetch('http://172.19.0.3:8080/check-location', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ latitude, longitude, human_sensor, light_sensor}),
    })
    .then(response => {
        if(response.ok){
            document.getElementById('status').textContent = "位置情報を更新しました。"
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('status').textContent = "判定に失敗しました。";
    });
}

function putUserLocation(position) {
    const latitude = position.coords.latitude;
    const longitude = position.coords.longitude;
    document.getElementById('status').textContent = "位置情報を送信中...";

    fetch('http://172.19.0.3:8080/put-user-location', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ latitude, longitude }),
    })
    .then(response => {
        if(response.ok){
            document.getElementById('status').textContent = "位置情報を更新しました。"
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('status').textContent = "判定に失敗しました。";
    });
}

function putHomeLocation(position) {
    const latitude = position.coords.latitude;
    const longitude = position.coords.longitude;
    document.getElementById('status').textContent = "位置情報を送信中...";

    fetch('http://172.19.0.3:8080/put-home-location', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ latitude, longitude }),
    })
    .then(response => {
        if(response.ok){
            document.getElementById('status').textContent = "位置情報を更新しました。"
        }
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('status').textContent = "判定に失敗しました。";
    });
}

function error() {
    document.getElementById('status').textContent = "位置情報の取得に失敗しました。";
}

document.getElementById('changeEmailBtn').addEventListener('click', () => {
    document.getElementById('changeEmailForm').style.display = "block";
});

document.getElementById('submitEmailBtn').addEventListener('click', () => {
    const newEmail = document.getElementById('newEmail').value;

    fetch('http://172.19.0.3:8080/change-email', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': localStorage.getItem('token'),
        },
        body: JSON.stringify({email:newEmail}),
    })
    .then(response => response.json())
    .then(data => {
        console.log(data);
        document.getElementById('email-change-message').textContent = "メールアドレスを変更しました。";
        localStorage.removeItem('token');
        document.getElementById('auth-container').style.display = "block";
        document.getElementById('main-container').style.display = "none";
        
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('email-change-message').textContent = "メールアドレス変更に失敗しました。";
    });
})
document.getElementById('logoutBtn').addEventListener('click', () => {
    localStorage.removeItem('token');
    window.location.reload();
})