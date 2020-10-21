
class Avatar {
    constructor(x, y, angle) {
        this.x = x;
        this.y = y;
        this.angle = angle;
        const boardX = this.x;
        const boardY = this.y + 50 + 20;
        this.board = new Board([{ value: 12, color: 1 }, { value: 7, color: 0 }]);
        this.board.setup(boardX, boardY);
    }

    draw() {
        fill(255)
        ellipse(this.x, this.y, 50, 50);
        this.board.draw();
    }
}