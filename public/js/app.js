const app = (() => {
    function applyDocBehavior(target) {
        jQuery('.list_remove', target).click(onListItemRemove);
        jQuery('.list_add', target).click(onListItemAdd);
        jQuery('.breakpoint_remove', target).click(onBreakpointRemove);
    }

    function onAppInit() {
        applyDocBehavior(document);
    }

    function onListItemRemove() {
        jQuery(this).parents('li').remove();
    }

    function onListItemAdd() {
        const link = jQuery(this);
        const name = link.attr('xattr_with_name');
        const list = jQuery('ul:first', link.parent('div'));
        const newItem = jQuery(newListItemString(name));

        applyDocBehavior(newItem);
        list.append(newItem);
    }

    function onBreakpointRemove() {
        jQuery(this).parents('.config').remove();
    }

    function newListItemString(name) {
        return `
            <li>
                <input type="text" name="${name}" />
                <a class="list_remove" href="#nogo">⛔</a>
            </li>
        `
    }

    return { onAppInit: onAppInit };
})();

jQuery(app.onAppInit);
