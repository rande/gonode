{% extends "dashboard:layouts/default.tpl" %}

{% block content %}
    <h1>Create a new node</h1>

    {% if error %} 
        <div class="error">{{ error }}</div>
    {% endif %}

    {% if !form %}
        <div>
            <ul>
                {% for key in keys %}
                    <li><a href="{{ url('dashboard_node_create', url_values('type', key)) }}">{{ key }}</a></li>
                {% endfor %}
            </ul>
        </div>
    {% endif %}


   {% if form %}
        <form 
            action="{{ url('dashboard_node_create') }}" 
            method='POST'>

            <div class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="name">
                    Name
                </label>

                <input 
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                    name="Name" 
                    type="text" 
                    value="{{form.Name}}"
                    placeholder="Name"
                >

                {% if (errors && errors.Name)%}
                    <p class="text-red-500 text-xs italic">Error: {{ errors.Name.Error() }}</p>
                {% endif %}

                <label class="block text-gray-700 text-sm font-bold mb-2" for="slug">
                    Slug
                </label>

                <input 
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                    name="Slug" 
                    type="text" 
                    value="{{form.Slug}}"
                    placeholder="slug-for-the-url"
                >
                    
                {% if (errors && errors.Slug)%}
                    <p class="text-red-500 text-xs italic"> Error: {{ errors.Slug.Error() }}</p>
                {% endif %}

                <label class="block text-gray-700 text-sm font-bold mb-2" for="status">
                    Status
                </label>

                <input 
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                    name="Status"
                    type="text" 
                    value="{{form.Status}}"
                >
                {% if (errors && errors.Status)%}
                    <p class="text-red-500 text-xs italic"> Error: {{ errors.Status.Error() }}</p>
                {% endif %}

                <label class="block text-gray-700 text-sm font-bold mb-2" for="weight">
                    Weight
                </label>

                <input 
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                    name="Weight" 
                    type="text" 
                    value="{{form.Weight}}"
                >
                {% if (errors && errors.Weight)%}
                    <p class="text-red-500 text-xs italic"> Error: {{ errors.Weight.Error() }}</p>
                {% endif %}

                <label class="block text-gray-700 text-sm font-bold mb-2" for="enabled">
                    Enabled
                </label>

                <input 
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                    name="Enabled" 
                    type="text" 
                    value="{{form.Enabled}}"
                >
                {% if (errors && errors.Enabled)%}
                    <p class="text-red-500 text-xs italic"> Error: {{ errors.Enabled.Error() }}</p>
                {% endif %}

                <label class="block text-gray-700 text-sm font-bold mb-2" for="ParentUuid">
                    ParentUuid
                </label>

                <input 
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                    name="ParentUuid" 
                    type="text" 
                    value="{{form.ParentUuid}}"
                >

                {% if (errors && errors.ParentUuid)%}
                    <p class="text-red-500 text-xs italic"> Error: {{ errors.ParentUuid.Error() }}</p>
                {% endif %}

                <div class="flex items-center justify-between">
                    <input 
                        type='submit' 
                        value='Create' 
                        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" 
                    />
                </div>
            </div>
        </form>
    {% endif %}
{% endblock %}