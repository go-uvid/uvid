const offscreen = new OffscreenCanvas(screen.width, screen.height);
const context = offscreen.getContext(
	'2d',
) as unknown as CanvasRenderingContext2D;

export function measureTextWidth(text: string, font = '') {
	context?.save();
	context.font = font;
	context?.restore();
	return context?.measureText(text).width;
}
