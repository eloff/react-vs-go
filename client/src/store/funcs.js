export const dateFromSlot = (slot) => {
	const args = slot.Date.split('-').map(x => parseInt(x, 10));
	args[1]--; // javascript months are zero-based, but go months are one-based
	args.push(Math.floor(slot.StartMinute / 60));
	args.push(slot.StartMinute % 60);
	return new Date(...args);
};