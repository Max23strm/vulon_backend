package queries

var EventsByMonth = `SELECT 
	event_uid, 
	event_type, 
	date_start, 
	date_end, 
	events.time_zone,
	event_name, 
	event_types.type_name, 
	events_state.event_state 
	FROM events 
	INNER JOIN event_types ON events.event_type = event_types.event_type_uid 
	INNER JOIN events_state ON events.event_state_uid = events_state.event_state_uid 
	WHERE entity_uid = $1  
	AND events.date_start >= $2 
	AND events.date_start < $3;`
