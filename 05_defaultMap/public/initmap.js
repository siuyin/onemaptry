let sw = L.latLng(1.144, 103.535);
let ne = L.latLng(1.494, 104.502);
let bounds = L.latLngBounds(sw, ne);


let ctr = L.latLng(1.3461151,103.773155);
var map = L.map('mapdiv', {
  center: ctr,
  zoom: 17
});

map.setMaxBounds(bounds);

let basemap = L.tileLayer('https://www.onemap.gov.sg/maps/tiles/Default_HD/{z}/{x}/{y}.png', {
  detectRetina: true,
  maxZoom: 19,
  minZoom: 11,
});

basemap.addTo(map);
L.marker(ctr)
  .bindPopup("Lobby B<br>31 Hindhede Walk<br>Southaven 2")
  .addTo(map)
//.openPopup()

function onMapClick(e) {
  let mypop=L.popup()
  mypop
    .setLatLng(e.latlng)
    .setContent(`${e.latlng.toString()}`)
    .openOn(map)

}

//map.on("click",onMapClick)
var redicon = L.icon({
  iconUrl: 'redmarker.png',
  iconSize: [26,41],
  iconAnchor: [13,40],
  popupAnchor: [0,-25]
})

function placeRedicon(e) {
  L.marker(e.latlng,{icon: redicon})
    .bindPopup(`${e.latlng.toString()}`)
    .addTo(map)
}
//map.on("click",placeRedicon)
