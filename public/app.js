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
        data.forEach(function(p, i) {
            var html = '<div class="panel panel-default">' +
                '  <div class="panel-heading">' + '<h3 class="panel-title">' + p.title + '</h3>' + '</div>' +
                '  <div class="panel-body">' + p.body + '</div>' +
                '  <div class="panel-footer">' +
                '    <small>Posted: ' + p.time + '</small>' +
                '    <button class="del btn btn-danger pull-right" data-id="' + i + '">Delete</button>' +
                '    <div class="clearfix"></div>' +
                '  </div>' +
                '</div>';

            $("#post-list").prepend(html);
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
            body: $("#body").val(),
            title: $("#title").val()
        };

        $.ajax({
            type: 'POST',
            url: '/posts',
            data: JSON.stringify(data),
        })
        .then(function() {
            $("#new-post")[0].reset();
            updateList();
        })
        .fail(function() {
            alert("Failed to add Post.");
        });
    });

    $("#post-list").click(function(e){
        var target = $(e.originalEvent.target);
        if (!target.hasClass("del")) {
            return;
        }

        var id = $(target).data("id");
        console.log("Deleting post", id);

        $.ajax({
            type: 'DELETE',
            url: '/posts/'+id
        })
        .then(function() {
            updateList();
        })
        .fail(function() {
            alert("Failed to delete Post.");
        });
    });

    // Make the initial call to update our list of posts
    updateList();
});
