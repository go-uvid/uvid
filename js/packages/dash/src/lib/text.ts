const offscreen = new OffscreenCanvas(screen.width, screen.height);
const context = offscreen.getContext('2d');

export function measureTextWidth(text: string, font = '') {
	context?.save();
	context!.font = font;
	return context?.measureText(text).width;
}
