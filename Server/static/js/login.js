// ---------- AUTO REDIRECT IF ALREADY LOGGED IN ----------
const token = localStorage.getItem('jwt');
if(token && window.location.pathname.includes('login.html')){
    window.location.href = "index.html";
}

const loginForm = document.getElementById('loginForm');
const registerForm = document.getElementById('registerForm');

const API = "http://localhost:8082";

// ---------- LOGIN ----------
loginForm.addEventListener('submit', async (e)=>{
    e.preventDefault();
    const username = document.getElementById('loginUsername').value;
    const password = document.getElementById('loginPassword').value;

    try {
        const res = await fetch(`${API}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });

        if(res.ok){
            const data = await res.json();
            const token = data.token || data.accessToken;
            if(!token) throw new Error("No token received from backend");

            localStorage.setItem('jwt', token);

            // Redirect instantly
            window.location.href = "/dashboard/dashboard.html";
        } else {
            let errMsg = res.statusText;
            try {
                const data = await res.json();
                if(data.message) errMsg = data.message;
            } catch(e){}
            alert("Login failed: " + errMsg);
        }

    } catch(err){
        alert("Login error: " + err.message);
    }
});




// ---------- REGISTER ----------
registerForm.addEventListener('submit', async (e)=>{
    e.preventDefault();
    const username = document.getElementById('regUsername').value;
    const password = document.getElementById('regPassword').value;

    try{
        const res = await fetch(`${API}/registration`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });

        if(res.ok){
            alert("Registration successful! You can now login.");
            registerForm.reset();
        } else {
            let errMsg = res.statusText;
            try{
                const data = await res.json();
                if(data.message) errMsg = data.message;
            } catch(e){}
            alert("Registration failed: " + errMsg);
        }

    } catch(err){
        alert("Registration error: " + err.message);
    }
});
