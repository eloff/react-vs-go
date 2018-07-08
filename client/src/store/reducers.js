const initialState = {
	"ClientID":2,
	"Company":{"ID":1},
	"Days": []
};

const rootReducer = (state = initialState, action) => {
	if (action && action.data) {
		return action.data;
	}
	return state;
};

export default rootReducer;