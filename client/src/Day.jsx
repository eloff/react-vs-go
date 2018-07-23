import React, { Component } from 'react';
import { compose } from 'redux';
import { connect } from 'react-redux';
import { makeReservation } from './store/reserve';
import { dateFromSlot } from './store/funcs';
import PerformanceHoc from './PerformanceHoc'

const displayTime = (slot) => {
	const hour = Math.floor(slot.StartMinute / 60);
	const minute = slot.StartMinute % 60;
	return hour.toString().padStart(2, '0') + ':' + minute.toString().padStart(2, '0');
};

const Day = ({ clientID, slots, onStartPerformanceTimer, makeReservation }) => {
    const date = dateFromSlot(slots[0]);
	const now = new Date();
	const onReserve = reservationDate => {
		onStartPerformanceTimer('Reservation');
		makeReservation(reservationDate);
	};

    return (
        <div className="day">
            <h4>{date.toLocaleDateString("en-us", {weekday: 'short', month: 'short', day: 'numeric'})}</h4>
            {slots.map((_, i) => i).filter(i => i % 2 === 0).map(index => (
				<div className="hour">
					<Slot clientID={clientID} now={now} onReserve={onReserve} slot={slots[index]} />
					<Slot clientID={clientID} now={now} onReserve={onReserve} slot={slots[index + 1]} />
				</div>
			))}
        </div>
    );
};

const Slot = ({ clientID, now, slot, onReserve }) => {
	const date = dateFromSlot(slot);

	if (now > date) {
		return <div className="reserved past">00:00</div>;
	}

	if (slot.Available) {
		return <button onClick={() => onReserve(date)}>{displayTime(slot)}</button>;
	}

	return <div className={(slot.ClientID === clientID) ? 'mine' : 'reserved'}>{displayTime(slot)}</div>;
};

const mapStateToProps = state => ({ clientID: state.ClientID });
const mapDispatchToProps = { makeReservation };
const areOwnPropsEqual = (oldProps, newProps) =>
	Object.keys(oldProps.slots).find(i => oldProps.slots[i].Available !== newProps.slots[i].Available) === undefined;

export default compose(
	connect(mapStateToProps, mapDispatchToProps, undefined, { areOwnPropsEqual }),
	PerformanceHoc,
)(Day);
