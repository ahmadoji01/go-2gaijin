/*
// topbar dropdown
*/
$('.topbar-dropdown__btn').on('click', function() {
$(this).closest('.topbar-dropdown').toggleClass('topbar-dropdown--opened');
});

document.addEventListener('click', function(event) {
$('.topbar-dropdown')
.not($(event.target).closest('.topbar-dropdown'))
.removeClass('topbar-dropdown--opened');
}, true);

touchClick(document, function(event) {
$('.topbar-dropdown')
.not($(event.target).closest('.topbar-dropdown'))
.removeClass('topbar-dropdown--opened');
});
