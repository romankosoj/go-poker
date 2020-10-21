let items = [];
let angles = [];
let printed = false;

let board;

let tableW;
let tableH;

function setup() {
    createCanvas(window.innerWidth, window.innerHeight - 50)



    tableW = width - 1000;
    tableH = height - 450;
    generateItems(10, tableW / 2, tableH / 2)
    board = new Board(
        [
            { value: 10, color: 3 },
            { value: 0, color: 1 },
            { value: 11, color: 0 },
            { value: 11, color: 0 },
            { value: 11, color: 0 },
        ]);
    //console.log("BoardW", (width / 2) - (board.w / 2))
    board.setup((width / 2) - (board.w / 2), (height / 2) - (board.h / 2))
}


// (x, y) = (rx * cos(θ), ry * sin(θ))
// (x, y) = ((width - 500) * cos(θ), (height - 200) * sin(θ))
function draw() {
    background(220)

    stroke(0, 0);
    fill(81, 200, 45)
    ellipse(width / 2, height / 2, width - 1000, height - 450);

    for (let i = 0; i < items.length; i++) {
        items[i].draw();
    }
    board.draw();


}

function generateItems(n, rx, ry) {
    angles = [];
    items = [];
    var frags = 360 / n;
    for (let i = 0; i <= n; i++) {
        angles.push(frags * i * Math.PI / 180);
        //angles.push(frags * i * Math.PI);
    }

    for (let i = 0; i < n; i++) {
        const x = (width / 2) + rx * Math.cos(angles[i])
        const y = (height / 2) + ry * Math.sin(angles[i])
        let pl = new Avatar(x, y, angles[i])
        items.push(pl)
    }

}

function degreeToRadians($degree) {
    return $degree * Math.PI / 180;
}