function showMessage(type, summary, content) {
    document.getElementById('message_combo').setAttribute("class", "message "+type);
    document.getElementById('message_summary').textContent = summary;
    document.getElementById('message_content').textContent = content;
}

function clearMessage() {
    document.getElementById('message_combo').setAttribute("class", "message");
    document.getElementById('message_summary').textContent = "";
    document.getElementById('message_content').textContent = "";
}