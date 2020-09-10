play.onclick = ev => {
  if (document.querySelectorAll('.player-name').length < 2) {
    alert('Not enough players');
    ev.preventDefault();
  }
};

const connect = () => {
    const url = location.origin.replace('http', 'ws') + webSocketURLH.value;
    const socket = new WebSocket(url);

    socket.onmessage = ev => {
        console.log(ev);
        location.reload();
    };

    socket.onerror = ev => {
        setTimeout(connect, 1000);
    };
};

connect();
