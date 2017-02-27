/**
 * Created by tuanpa on 1/31/17.
 */






// jQuery(document).ready(function(){
//     single_post_md_parse(); // call this before loadTorrent
//     loadTorrent();
// });

///////////////////////////////
// WebTorrent Video
var webTorrent;
var counter = 0;

function initWebTorrent() {
    if (typeof variable === 'undefined') {
        webTorrent = new WebTorrent();
    }
}

function loadTorrent() {
    jQuery("img").each(function( index ) {
        var img = jQuery( this );
        // var torrent_alt = img.attr('alt');
        var torrent_uri = img.attr('src');


        //if (!torrent_uri.match("^/torrent/")) { // dont use magnet anymore, now use only torrent
        if (!torrent_uri.match(".torrent$")) { // check end with .torrent
            return true; // continue next loop
        }

        torrent_uri = "/torrent/" + torrent_uri;
        torrent_uri = absolutePath(torrent_uri);
        console.log( "torrent_uri=" + torrent_uri );

        // var torrent_obj = parseQuery(torrent_uri);
        // console.log(torrent_obj)
        // var file_ext = magnet_obj["type"].split('.').pop();
        // var file_ext = torrent_obj["type"];
        // console.log( "file_ext=" + file_ext );

        // video (.mp4, .webm, .m4v, etc.),
        // audio (.m4a, .mp3, .wav, etc.)
        // images (.jpg, .gif, .png, etc.),
        // and other file formats (.pdf, .md, .txt, etc.).
        // if (jQuery.inArray(file_ext, ["mp4", "avi", "mov", "wmv", "mpeg"]) != -1) {
            //console.log("video");
            loadTorrentVideo (img, torrent_uri);
            counter++;
        // }
    });
}

function loadTorrentVideo (img, torrent_uri) {
    initWebTorrent();

    var outputId = "output" + counter;
    var progressBarId = "progressBar" + counter;

    var root = jQuery('<div id="' + outputId + '"><div id="' + progressBarId + '"></div></div>');
    img.replaceWith( root );

    webTorrent.add(torrent_uri, function (torrent) {
        console.log("torrent loaded");
        var file = torrent.files[0];   // Torrents can contain many files. Let's use the first.
        file.appendTo( "#" + outputId ); // Display the file by adding it to the DOM. Supports video, audio, image, etc. files
    });
}

function parseQuery(qstr) {
    var query = {};
    var a = (qstr[0] === '?' ? qstr.substr(1) : qstr).split('&');
    for (var i = 0; i < a.length; i++) {
        var b = a[i].split('=');
        query[decodeURIComponent(b[0])] = decodeURIComponent(b[1] || '');
    }
    return query;
}

/**
 * Relative URL to absolute URL
 * @param href
 * @returns {string}
 */
function absolutePath(href) {
    var link = document.createElement("a");
    link.href = href;
    return (link.protocol+"//"+link.host+link.pathname+link.search+link.hash);
}

///////////////
// for single post page
function single_post_md_parse(Config_IpFsGateway) {
    // post
    if (jQuery("#PostContent").length <= 0) {
        return;
    }

    // http://127.0.0.1:8080/ipfs/QmR168W81xikZdtCfzYYp7TVi1L9Cad2UTbLmTPt3e73PT
    var url = Config_IpFsGateway + "/ipfs/" + jQuery("#PostContent").text();

    $.get( url, function( data ) {
        var md_str = data;
        var html_str = SimpleMDE.prototype.markdown(md_str);
        jQuery("#PostContent").parent().html(html_str);

        // comments
        jQuery( ".comment-container > textarea" ).each(function() {
            var comment_md_str = jQuery( this ).text();
            var comment_html_str = SimpleMDE.prototype.markdown(comment_md_str);
            console.log(comment_html_str);
            jQuery(this).parent().html(comment_html_str);
        });

        loadTorrent();
    });

}