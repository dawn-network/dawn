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
        var magnet_uri = img.attr('src');
        //console.log( magnet_uri );

        if (!magnet_uri.match("^magnet")) { // begins with magnet
            return true; // continue next loop
        }

        var magnet_obj = parseQuery(magnet_uri)
        var file_ext = magnet_obj["dn"].split('.').pop();

        // video (.mp4, .webm, .m4v, etc.),
        // audio (.m4a, .mp3, .wav, etc.)
        // images (.jpg, .gif, .png, etc.),
        // and other file formats (.pdf, .md, .txt, etc.).
        if (jQuery.inArray(file_ext, ["mp4", "avi", "mov", "wmv", "mpeg"]) != -1) {
            //console.log("video");
            loadTorrentVideo (img, magnet_uri);
            counter++;
        }
    });
}

function loadTorrentVideo (img, magnet_uri) {
    initWebTorrent();

    var outputId = "output" + counter;
    var progressBarId = "progressBar" + counter;

    var root = jQuery('<div id="' + outputId + '"><div id="' + progressBarId + '"></div></div>');
    img.replaceWith( root );

    webTorrent.add(magnet_uri, function (torrent) {
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



///////////////
// for single post page
function single_post_md_parse() {
    // post
    if (jQuery("#PostContent").length <= 0) {
        return;
    }

    var md_str = jQuery("#PostContent").text();
    var html_str = SimpleMDE.prototype.markdown(md_str);
    jQuery("#PostContent").parent().html(html_str);

    // comments
    jQuery( ".comment-container > textarea" ).each(function() {
        var comment_md_str = jQuery( this ).text();
        var comment_html_str = SimpleMDE.prototype.markdown(comment_md_str);
        console.log(comment_html_str);
        jQuery(this).parent().html(comment_html_str);
    });
}