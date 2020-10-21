const cardWidth = 75;
const cardHeight = 110;

class Card {
    constructor(value, color) {
        this.v = value;
        this.c = color
    }

    setup() {
        this.img = loadImage(`/public/cards/${this.v}_${this.c}.svg`);
    }

    draw(x, y) {
        image(this.img, x, y, cardWidth, cardHeight);
    }

}