
// PlayButton just draws a play button an thumbnail for the image before the video has played
export class PlayButton {
    
    private hasPlayed = false;
    private image = new Image();

    constructor(
        videoID: number,
        private ctx: CanvasRenderingContext2D,
        private width: number,
        private height: number
    ) {
        this.image.src = `/app/video/thumbnail?id=${videoID}`;
        this.image.onload = () => {
          this.draw();
        };
        this.draw();
    }

    private draw() {
        if (this.hasPlayed) {
            return;
        }
        this.ctx.clearRect(0, 0, this.width, this.height);
        if (this.image.complete) {
            let t = Math.max(this.height - (this.image.height * this.width) / this.image.width, 0) / 2;
            let l = Math.max(this.width - (this.image.width * this.height) / this.image.height, 0) / 2;
            this.ctx.drawImage(this.image, l, t, this.width - l, this.height - t);
        }
        const minDim = Math.min(this.width, this.height);
        const cx = this.width / 2;
        const cy = this.height / 2;
        this.ctx.beginPath();
        this.ctx.arc(cx, cy, minDim * .4, 0, 6.2832, false);

        this.ctx.strokeStyle = "black";
        this.ctx.lineWidth = minDim * 0.06;
        this.ctx.stroke();
        this.ctx.strokeStyle = "white";
        this.ctx.lineWidth = minDim * 0.05;
        this.ctx.stroke();

        this.ctx.beginPath();
        this.ctx.moveTo(cx - minDim * .18, cy - minDim * .2);
        this.ctx.lineTo(cx + minDim * .22, cy);
        this.ctx.lineTo(cx - minDim * .18, cy + minDim * .2);
        this.ctx.closePath();

        this.ctx.lineWidth = minDim * 0.01;
        this.ctx.strokeStyle = "black";
        this.ctx.stroke();
        this.ctx.fillStyle = 'white';
        this.ctx.fill();
    }

    public cancel() {
        this.hasPlayed = true;
    }
}