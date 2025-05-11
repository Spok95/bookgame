// luck.js — бросок кубика с анимацией и результатом удачи

const diceImages = [
    '/static/images/spin-1.png', '/static/images/spin-2.png', '/static/images/spin-3.png',
    '/static/images/spin-4.png', '/static/images/spin-5.png', '/static/images/spin-6.png',
    '/static/images/spin-7.png', '/static/images/spin-8.png', '/static/images/spin-9.png',
    '/static/images/spin-10.png', '/static/images/spin-11.png', '/static/images/spin-12.png',
    '/static/images/spin-13.png', '/static/images/spin-14.png', '/static/images/spin-15.png',
    '/static/images/spin-16.png'
];

const resultImages = [
    '/static/images/dice-1.png',
    '/static/images/dice-2.png',
    '/static/images/dice-3.png',
    '/static/images/dice-4.png',
    '/static/images/dice-5.png',
    '/static/images/dice-6.png'
];

async function checkLuck() {
    const button = document.querySelector("#luck-section button");
    const messageBox = document.getElementById("luck-message");
    const linksBox = document.getElementById("luck-links");
    const successLink = document.getElementById("luck-success");
    const failLink = document.getElementById("luck-fail");

    button.disabled = true;
    messageBox.textContent = "Бросаем кубик...";

    const diceImg = document.createElement("img");
    diceImg.id = "luck-dice";
    diceImg.src = "/static/images/dice-1.png";
    diceImg.width = 64;
    diceImg.style.display = "block";
    diceImg.style.margin = "1em auto";
    messageBox.innerHTML = "";
    messageBox.appendChild(diceImg);

    // 🎲 Анимация броска
    let frame = 0;
    let delay = 30;
    for (let i = 0; i < 20; i++) {
        frame = (frame + 1) % diceImages.length;
        diceImg.src = diceImages[frame];
        await new Promise(resolve => setTimeout(resolve, delay));
        delay += 5;
    }

    // Запрос к серверу
    const success = successLink.getAttribute("href");
    const fail = failLink.getAttribute("href");

    fetch(`/luck?success=${encodeURIComponent(success)}&fail=${encodeURIComponent(fail)}`, {
        method: "POST"
    })
        .then(res => res.json())
        .then(data => {
            diceImg.src = resultImages[data.Roll - 1];
            messageBox.innerHTML += `<p>Выпало: <strong>${data.Roll}</strong></p>`;

            if (data.lucky) {
                messageBox.innerHTML += `<p>✅ ${data.message}</p>`;
                successLink.style.display = "inline-block";
                failLink.style.display = "none";
            } else {
                messageBox.innerHTML += `<p>❌ ${data.message}</p>`;
                successLink.style.display = "none";
                failLink.style.display = "inline-block";
            }

            linksBox.style.display = "block";
        })
        .catch(err => {
            messageBox.textContent = "Ошибка проверки удачи.";
            console.error(err);
        });
}