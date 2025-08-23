const canvas = document.getElementById('sakuraCanvas');
const ctx = canvas.getContext('2d');
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

const petals = [];
for(let i=0;i<50;i++){
    petals.push({
        x: Math.random()*canvas.width,
        y: Math.random()*canvas.height,
        r: Math.random()*5+2,
        d: Math.random()*1+1,
        swing: Math.random()*2
    });
}

function drawPetals(){
    ctx.clearRect(0,0,canvas.width,canvas.height);
    ctx.fillStyle = "rgba(255,182,193,0.8)";
    ctx.beginPath();
    petals.forEach(p=>{
        ctx.moveTo(p.x,p.y);
        ctx.arc(p.x,p.y,p.r,0,Math.PI*2);
    });
    ctx.fill();
    updatePetals();
}

function updatePetals(){
    petals.forEach(p=>{
        p.y += p.d;
        p.x += Math.sin(p.y*0.01)*p.swing;
        if(p.y>canvas.height){ p.y=-10; p.x=Math.random()*canvas.width; }
    });
    requestAnimationFrame(drawPetals);
}

drawPetals();
window.addEventListener('resize',()=>{
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
});
