const diceSound = new Audio('/static/music/dice-roll.mp3');
const dice = document.getElementById('dice');
const rollButton = document.getElementById('roll-button');
const battleLog = document.getElementById('battle-log');

const faces = [
    '/static/images/dice-1.png',
    '/static/images/dice-2.png',
    '/static/images/dice-3.png',
    '/static/images/dice-4.png',
    '/static/images/dice-5.png',
    '/static/images/dice-6.png'
];

const spins = [
    '/static/images/spin-1.png',
    '/static/images/spin-2.png',
    '/static/images/spin-3.png',
    '/static/images/spin-4.png',
    '/static/images/spin-5.png',
    '/static/images/spin-6.png',
    '/static/images/spin-7.png',
    '/static/images/spin-8.png',
    '/static/images/spin-9.png',
    '/static/images/spin-10.png',
    '/static/images/spin-11.png',
    '/static/images/spin-12.png',
    '/static/images/spin-13.png',
    '/static/images/spin-14.png',
    '/static/images/spin-15.png',
    '/static/images/spin-16.png'
];

async function startAttack() {
    if (!dice || !rollButton || !battleLog) return;

    rollButton.disabled = true;

    // 🔥 Проигрываем звук броска
    diceSound.currentTime = 0;
    diceSound.play();

    let frame = 0;
    let totalFrames = 20;
    let delay = 20;
    let bounceAmplitude = 20; // Начальная амплитуда прыжка (в пикселях)

    // Анимация вращения с замедлением + подпрыгивание
    for (let i = 0; i < totalFrames; i++) {
        frame = (frame + 1) % spins.length;
        dice.src = spins[frame];

        // Прыжок вверх и вниз
        const direction = (i % 2 === 0) ? 1 : -1; // чередуем вверх/вниз
        dice.style.transform = `translateY(${direction * bounceAmplitude}px)`;

        await new Promise(resolve => setTimeout(resolve, delay));

        delay += 5;            // плавное замедление
        bounceAmplitude *= 0.9; // амплитуда прыжка уменьшается каждый кадр
    }

    try {
        const response = await fetch('/attack', { method: 'POST' });
        const data = await response.json();

        const playerRoll = data.PlayerRoll;
        const enemyRoll = data.EnemyRoll;
        const resultText = data.Result;

        // Показать итоговую грань от сервера
        dice.src = faces[playerRoll - 1];
        dice.style.transform = 'translateY(0)';

        // Обновить лог боя
        const logEntry = document.createElement('p');
        logEntry.textContent = `Вы: ${playerRoll} | Враг: ${enemyRoll} → ${resultText}`;
        battleLog.appendChild(logEntry);
        battleLog.scrollTop = battleLog.scrollHeight;

        if (data.BattleEnded) {
            const endMsg = document.createElement('p');
            endMsg.innerHTML = `<strong>Бой окончен!</strong>`;
            battleLog.appendChild(endMsg);
            rollButton.disabled = true;
        } else {
            rollButton.disabled = false;
        }
    } catch (error) {
        console.error('Ошибка атаки:', error);
        rollButton.disabled = false;
    }
}


