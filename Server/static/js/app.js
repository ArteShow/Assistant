// ---------- DASHBOARD PROTECTION ----------
const jwtToken = localStorage.getItem('jwt');
if (!jwtToken) {
    alert("You must login first!");
    window.location.href = "login.html";
}

// ---------------- MONEY & TASKS ----------------
let moneyGoal = 13500;
let moneyCollected = 0;

const moneyBar = document.getElementById('moneyBar');
const moneyInfo = document.getElementById('moneyInfo');

let tasks = [];
const tasksContainer = document.getElementById('tasksContainer');

// ---------------- TASK RENDERING ----------------
function renderTasks(filter = 'all') {
    tasksContainer.innerHTML = '';
    tasks
        .filter(t => filter === 'all' ? true : t.status.toLowerCase() === filter)
        .forEach(task => {
            const div = document.createElement('div');
            div.className = 'task-card';
            div.innerHTML = `
                <span class="status-badge ${task.status.toLowerCase() === 'finished' ? 'status-finished' : 'status-unfinished'}">
                    ${task.status.toUpperCase()}
                </span>
                <p><strong>ID:</strong> ${task.id}</p>
                <p><strong>Title:</strong> ${task.title}</p>
                <p><strong>Description:</strong> ${task.description}</p>
                <p><strong>Worth:</strong> €${task.money}</p>
                <div class="task-buttons">
                    <button class="finish" onclick="finishTask(${task.id}, ${task.money})">Finish</button>
                    <button class="delete" onclick="deleteTask(${task.id})">Delete</button>
                </div>
            `;
            tasksContainer.appendChild(div);
        });
    updateMoney();
}

// ---------------- MONEY CHART ----------------
function updateMoney() {
    moneyInfo.textContent = `Collected: €${moneyCollected} / Goal: €${moneyGoal}`;
}

function animateMoneyBar() {
    const percent = moneyGoal > 0 ? Math.min((moneyCollected / moneyGoal) * 100, 100) : 0;
    moneyBar.style.width = '0%';
    if (percent === 0) return;
    let start = null;
    function step(timestamp) {
        if (!start) start = timestamp;
        let progress = timestamp - start;
        let width = Math.min(progress / 15, percent);
        moneyBar.style.width = width + '%';
        if (width < percent) requestAnimationFrame(step);
    }
    requestAnimationFrame(step);
}

async function getMoneyStats() {
    const token = localStorage.getItem("jwt");
    try {
        const res = await fetch("http://localhost:8082/money/getStats", {
            method: "GET",
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        });

        if (res.ok) {
            const data = await res.json();
            moneyGoal = data.goal || 13500;
            moneyCollected = data.current_money || 0;
            updateMoney();
            animateMoneyBar();
        } else {
            console.error("Failed to fetch money stats");
        }
    } catch (err) {
        console.error(err);
    }
}

// ---------------- FETCH TASKS ----------------
async function getTasks() {
    const token = localStorage.getItem("jwt");
    try {
        const res = await fetch("http://localhost:8082/task/getUsersTaskList", {
            method: "POST",
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            }
        });

        if (res.ok) {
            const data = await res.json();
            tasks = Array.isArray(data) ? data : ((data && data.tasks) || []);
            renderTasks();
        } else {
            console.error("Failed to fetch tasks");
        }
    } catch (err) {
        console.error(err);
    }
}

// ---------------- FINISH & DELETE TASK ----------------
async function finishTask(id, money) {
    const token = localStorage.getItem("jwt");
    try {
        const res = await fetch("http://localhost:8082/task/editTasksStatus", {
            method: "POST",
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ task_id: id, status: "Finished" })
        });

        if (res.ok) {
            // Add money to goal in DB
            await fetch("http://localhost:8082/money/addMoney", {
                method: "POST",
                headers: {
                    "Authorization": "Bearer " + token,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({ amount: money })
            });

            // Update local variable immediately
            moneyCollected += money;
            updateMoney();
            animateMoneyBar();

            await getTasks();
        } else {
            console.error("Failed to update task status");
        }
    } catch (err) {
        console.error(err);
    }
}

async function deleteTask(id) {
    const token = localStorage.getItem("jwt");
    try {
        const res = await fetch("http://localhost:8082/task/delete", {
            method: "POST",
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ task_ID: id })
        });
        if (res.ok) getTasks();
        else console.error("Failed to delete task");
    } catch (err) {
        console.error(err);
    }
}

// ---------------- FILTER ----------------
function filterTasks(type) {
    renderTasks(type);
}

// ---------------- CREATE TASK MODAL ----------------
const modal = document.getElementById("taskModal");
const btn = document.getElementById("createTaskBtn");
const span = document.getElementsByClassName("close")[0];

btn.onclick = () => {
    modal.style.opacity = "1";
    modal.style.pointerEvents = "auto";
}

span.onclick = () => {
    modal.style.opacity = "0";
    modal.style.pointerEvents = "none";
}

window.onclick = (e) => {
    if (e.target == modal) {
        modal.style.opacity = "0";
        modal.style.pointerEvents = "none";
    }
}

// ---------------- CREATE TASK FORM ----------------
const taskForm = document.getElementById("taskForm");
taskForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const title = document.getElementById("taskTitle").value;
    const description = document.getElementById("taskDesc").value;
    const money = parseInt(document.getElementById("taskWorth").value);
    const status = document.getElementById("taskStatus").value;
    const token = localStorage.getItem("jwt");

    try {
        const res = await fetch("http://localhost:8082/task/add", {
            method: "POST",
            headers: {
                "Authorization": "Bearer " + token,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ title, description, money, status })
        });

        if (res.ok) {
            modal.style.opacity = "0";
            modal.style.pointerEvents = "none";
            taskForm.reset();
            await getTasks();
        } else {
            const data = await res.json();
            alert("Error: " + (data.message || res.statusText));
        }
    } catch (err) {
        alert("Request failed: " + err.message);
    }
});

// ---------------- INITIAL LOAD ----------------
getTasks();
getMoneyStats();
animateMoneyBar();

// ---------------- SAKURA PETALS ----------------
const canvas = document.getElementById('sakuraCanvas');
const ctx = canvas.getContext('2d');
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

const petals = [];
for (let i = 0; i < 50; i++) {
    petals.push({
        x: Math.random() * canvas.width,
        y: Math.random() * canvas.height,
        r: Math.random() * 5 + 2,
        d: Math.random() * 1 + 1,
        swing: Math.random() * 2
    });
}

function drawPetals() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
    ctx.fillStyle = "rgba(255,182,193,0.8)";
    ctx.beginPath();
    petals.forEach(p => {
        ctx.moveTo(p.x, p.y);
        ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
    });
    ctx.fill();
    updatePetals();
}

function updatePetals() {
    petals.forEach(p => {
        p.y += p.d;
        p.x += Math.sin(p.y * 0.01) * p.swing;
        if (p.y > canvas.height) { p.y = -10; p.x = Math.random() * canvas.width; }
    });
    requestAnimationFrame(drawPetals);
}

drawPetals();
window.addEventListener('resize', () => {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
});
