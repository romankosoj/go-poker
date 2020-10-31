const { Text, Sprite, Graphics, Container, Texture } = require("pixi.js");

class Notification extends Container {
    constructor(appWidth, appHeight) {
        super()
        this.text = new Text("Hello", {
            fontSize: 42,
            fill: 0xffffff
        });

        this.bg = new Graphics();

        this.textBg = new Graphics();

        this.appWidth = appWidth
        this.appHeight = appHeight;


        this.addChild(this.bg, this.textBg, this.text);

        this.update(appWidth, appHeight)
    }

    update() {

        this.bg.clear();
        this.bg.beginFill(0x000000, 0.5)
        this.bg.drawRect(0, 0, this.appWidth, this.appHeight)
        this.bg.endFill();

        const paddingX = 50;
        const paddingY = 30;
        this.textBg.clear();
        this.textBg.beginFill(0x000000, 0.75)
        this.textBg.drawRoundedRect(
            (this.appWidth / 2) - (this.text.width / 2) - (paddingX / 2),
            (this.appHeight / 2) - (this.text.height / 2) - (paddingY / 2),
            this.text.width + paddingX,
            this.text.height + paddingY, 8
        )
        this.textBg.endFill();

        this.text.position.set((this.appWidth / 2) - (this.text.width / 2), (this.appHeight / 2) - (this.text.height / 2))

    }
}

export { Notification }