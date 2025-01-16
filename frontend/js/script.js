document.getElementById('getLocationBtn').addEventListener('click', () => {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(success, error);
    } else {
        document.getElementById('status').textContent = "このブラウザは位置情報をサポートしていません。";
    }
});

function success(position) {
    const latitude = position.coords.latitude;
    const longitude = position.coords.longitude;
    const request_body = JSON.stringify({ latitude, longitude });
    console.log(request_body);
    document.getElementById('status').textContent = "位置情報を送信中...";

    fetch('http://172.18.0.3:8080/check-location', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: request_body
    })
    .then(response => {
        console.log(response);
    })
    .then(data => {
        console.log(data);
    })
    .catch(error => {
        console.error('Error:', error);
        document.getElementById('status').textContent = "判定に失敗しました。";
    });
}

function error() {
    document.getElementById('status').textContent = "位置情報の取得に失敗しました。";
}