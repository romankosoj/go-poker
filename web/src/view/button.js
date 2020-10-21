const { Container, Graphics, Text, graphicsUtils } = require("pixi.js");

class Button extends Container {
    constructor(options) {
        super();
        this.options = {
            roundedRadius: 5,

            background: 0x000000,
            hoverBackground: 0x252525,
            hover: false,
            fontSize: 35,
            label: "Button",
            fill: 0xffffff,
            paddingX: 10,
            paddingY: 5,
        }

        this.label = new Text("heelllo", {
            fill: 0xffffff,
        });
        this.label.position.set(0, 0)
        this.background = new Graphics();
        this.addChild(this.background, this.label);
        this.interactive = true;
        this.buttonMode = true;
        this.on("click", (e) => {
            this.update({ background: 0x444444, hoverBackground: 0x444444 })
        })
        this.on("mouseover", e => {
            this.update({ hover: true })
        });
        this.on("mouseout", e => {
            this.update({ hover: false })
        }.bind(this));
        this.update(options);
    }

    update(options) {
        this.options = {
            ...this.options,
            ...options,
        };

        this.onResize();
    }

    onResize() {
        this.background.clear();
        if (this.options.hover) {
            this.background.beginFill(this.options.hoverBackground, 1)
        } else {
            this.background.beginFill(this.options.background, 1)
        }
        this.background.drawRoundedRect(0, 0, this.label.width + 2 * this.options.paddingX, this.label.height + 2 * this.options.paddingY, this.options.roundedRadius)
        this.background.endFill();

        this.label.position.set(this.options.paddingX, this.options.paddingY);
    }

}

export { Button }