import { rW, rH, isMobile } from "./utils";

const { Text, Sprite, Graphics, Container, Texture } = require("pixi.js");

export const CARDWIDTH = 60;
export const CARDHEIGHT = 95;
const emptyCardText = "";

class Card extends Container {
    constructor(id, card) {
        super();

        this.card = {
            value: -1,
            color: -1,
        }

        if (isMobile()) {
            this.sprite = new Text(emptyCardText, {
                fontSize: 10,
            })
            this.addChild(this.sprite);
        } else {
            this.sprite = new Sprite(id["back.png"]);
            this.addChild(this.sprite);
            let rect = new Graphics()
            rect.lineStyle(1, 0x000000, 1);
            rect.drawRoundedRect(0, 0, rW(CARDWIDTH), rH(CARDHEIGHT), 5);
            this.addChild(rect)
        }





        this.update(card);
    }

    update(card) {
        this.card = {
            ...this.card,
            ...card,
        };
        if (isMobile()) {
            let uni;
            switch (this.card.color) {
                case 0:
                    uni = "	♣"
                    break;
                case 1:
                    uni = "♦"
                    break;
                case 2:
                    uni = "♥"
                    break;
                case 3:
                    uni = "♠"
                    break;
            }
            if (this.card.value === -1 || this.card.color === -1) {
                this.sprite.text = emptyCardText;
            } else {
                this.sprite.text = this.card.value + uni;
            }
        } else {
            if (this.card.color < 0 || this.card.value < 0) {
                this.sprite.texture = Texture.from("back.png");
            } else {
                this.sprite.texture = Texture.from(`${this.card.value}_${this.card.color}.png`);
            }
            this.sprite.texture.baseTexture.scaleMode = PIXI.SCALE_MODES.NEAREST;
            this.sprite.width = rW(CARDWIDTH);
            this.sprite.height = rH(CARDHEIGHT);
        }

    }
}

export { Card }