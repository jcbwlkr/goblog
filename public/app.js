$.ajaxSetup({
    dataType: 'json',
    contentType: 'application/json;charset=UTF-8'
});

function updateList() {
    $("#post-list").empty();
    $.ajax({
            type: 'GET',
            url: '/posts'
        })
            .then(function(data) {
                data.forEach(function(p) {
                    var html = "<div class='post well'>" + p.body + "<small class='time'>" + p.time + "</small></div>";
                    $("#post-list").append(html);
                });
            })
            .fail(function() {
                alert("Could not get posts!");
            });
}

$(function() {
    $("form").submit(function(e){
        e.preventDefault();

        var data = {
            body: $("textarea").val()
        };

        $.ajax({
            type: 'POST',
            url: '/posts',
            data: JSON.stringify(data),
        })
            .then(function() {
                updateList();
            })
            .fail(function() {
                alert("Failed to add Post.");
            });
    });

    // Make the initial call to update our list of posts
    updateList();
});