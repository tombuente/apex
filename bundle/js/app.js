import '@tabler/core'

const logisticsAddressDetailMap = (id, lat, lon) => {
	var map = L.map(id, {
		center: [lat, lon],
		zoom: 32
	});

	L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
		maxZoom: 19,
		attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
	}).addTo(map);

	L.marker([lat, lon]).addTo(map);
}

window.logisticsAddressDetailMap = logisticsAddressDetailMap;
