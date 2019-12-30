const app = (() => {
    function applyDocBehavior(target) {
        jQuery('.list_remove', target).on('click', app.onListItemRemove);
        jQuery('.list_add', target).on('click', app.onListItemAdd);
    }

    function onAppInit() {
        applyDocBehavior(document);
    }

    function onListItemRemove() {
        jQuery(this).parent('li').remove();
    }

    function onListItemAdd() {
        const link = jQuery(this);
        const name = link.attr('xattr_with_name');
        const list = jQuery('ul:first', link.parent('div'));
        const newItem = jQuery(newListItemString(name));

        applyDocBehavior(newItem);
        list.append(newItem);
    }

    function newListItemString(name) {
        return `
            <li>
                <input type="text" name="${name}" />
                <a class="list_remove" href="#nogo">[remove]</a>
            </li>
        `
    }

    return {
        onAppInit: onAppInit,
        onListItemRemove: onListItemRemove,
        onListItemAdd: onListItemAdd
    };
})();

jQuery(app.onAppInit);
