const { Container, Text, Graphics } = require("pixi.js");

class Notification extends Container {
    constructor(gameState, appWidth, appHeight) {
        super()

        this.notification =
            this.appWidth = appWidth
        this.appHeight = appHeight

        this.text = new Text("", { fill: 0xffffff });
        this.bg = new Graphics();
        this.textBg = new Graphics();

        this.addChild(this.bg)
        this.visible = false;
        gameState.setOnNotification(this.onNotification.bind(this))
    }

    update() {

        this.text.text = this.notification

        this.bg.clear();
        this.bg.beginFill(0x000000, 0.5);
        this.bg.drawRect(0, 0, this.appWidth, this.appHeight);
        this.bg.endFill();

        const paddingX = 50;
        const paddingY = 30;
        this.textBg.clear();
        this.textBg.beginFill(0x000000, 0.75);
        this.textBg.drawRoundedRect(
            (this.appWidth / 2) - (this.text.width / 2) - (paddingX / 2),
            (this.appHeight / 2) - (this.text.height / 2) - (paddingY / 2),
            this.text.width + paddingX,
            this.text.height + paddingY,
        )
        this.textBg.endFill();

        this.text.position.set(
            (this.appWidth / 2) - (this.text.width / 2),
            (this.appHeight / 2) - (this.text.height / 2),
        )
    }

    updateFromState() {
        this.update();
    }

    onNotification(text, st) {
        this.notification = text;
        this.update();
        this.visible = true;
        if (!st) {
            setTimeout(() => {
                this.visible = false;
            }, 2500);
        }
    }

    reset() {
        this.text = "";
        this.visible = false;
    }
}

export { Notification }