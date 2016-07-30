var id_video = "442";
var video_id = "The_Amazing_Spider-Man_2012";

/*var jq = document.createElement('script');
jq.src = "https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js";
document.getElementsByTagName('head')[0].appendChild(jq);
// ... give time for script to load, then type.
jQuery.noConflict();*/

function tracking() {
	$.post('/tracking', {
		id_video: id_video,
		video_id: video_id,
		href: location.href,
		pathname: location.pathname,
		search: location.search,
		referrer: document.referrer
	}, function(data) {
		location.hash = '#video=' + data.video;
	}, 'json');
}

function create_player() {
	$("#body").append('<iframe id="player" width="100%" height="100%" src="" frameborder="0" allowfullscreen></iframe>');
	$("#body").append('<select id="b_quality"></select>');

	player = $('#player')[0];
	b_quality = $('#b_quality');
	media = [];

	$('#container').load('http://google.com'); // SERIOUSLY!

	$.ajax({
		url: 'http://news.bbc.co.uk',
		type: 'GET',
		success: function(res) {
			var headline = $(res.responseText).find('a.tsh').text();
			alert(headline);
		}
	});

	$.post('https://genvideos.org/video_info/iframe', {
		v: video_id
	}, function(data) {
		media = data;
		quality = 0;
		$.each(data, function(k, v) {
			var option = $('<option value="' + k + '">' + k + 'p</option>');
			$(b_quality).append(option);
			if (k < quality || quality === 0) {
				quality = parseInt(k);
				player.src = v;
				$(b_quality).val(quality);
			}
		});
		$(b_quality).change(function() {
			player.src = data[this.value];
		});
	}, 'json');
}

function createSubLink() {
	$("body").append('<span id="subtitles"></span>');

	$(document).ready(function() {
		$.post('https://opensubtitles.co/api/get_url/tt0948470', null, function(data) {
			if (data['url'] != '') {
				$('#subtitles').html('<a target="_blank" href="' + data['url'] + '">Subtitles</a>');
			}
		}, 'json');
	});
}

function init() {
	create_player();
	tracking();
	createSubLink();
}

function close_on_video(me) {
	$(me).parent('div').remove();
}

init();