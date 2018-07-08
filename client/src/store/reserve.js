import fetch from 'cross-fetch';
import store from './index';
import { dateFromSlot } from "./funcs";

const SERVER = 'http://35.192.118.156:8080';
const RESERVATION_ENDPOINT = '/schedule/reserve';

function requestReservation() {
	return {
		type: 'REQUEST_RESERVATION'
	};
}

function madeReservation(state, date) {
	for (let i=0; i < state.Days.length; i++) {
		const day = state.Days[i];
		if (dateFromSlot(day[0]).getDate() !== date.getDate()) {
			continue;
		}
		for (let j=0; j < day.length; j++) {
			if (dateFromSlot(day[j]).getTime() === date.getTime()) {
				// Deep copy state with a copy of day, and then modify day
				// to belong to this ClientID and to no longer be available.
				const newState = {
					...state
				}
				newState.Days = state.Days.slice();
				newState.Days[i] = newState.Days[i].slice();
				newState.Days[i][j] = {
					...day[j],
					ClientID: state.ClientID,
					Available: false
				}
				state = newState;
				break;
			}
		}
	}
	return {
		type: 'MADE_RESERVATION',
		data: state
	};
}

export const makeReservation = (state, date) => {
	return (dispatch) => {
		dispatch(requestReservation());

		const companyID = state.Company.ID;
		const clientID = state.ClientID;
		const year = date.getFullYear();
		const month = date.getMonth() + 1; // javascript is zero-based, go is one-based
		const day = date.getDate();
		const hour = date.getHours();
		const minute = date.getMinutes();
		const url = (
			SERVER +
			RESERVATION_ENDPOINT +
			'?companyID=' + companyID +
			'&clientID=' + clientID +
			'&year=' + year +
			'&month=' + month +
			'&day=' + day +
			'&hour=' + hour +
			'&minute=' + minute
		);
		return fetch(url).then(
			res => dispatch(madeReservation(state, date)),
			err => console.error('error in makeReservation', err)
		);
	};
};