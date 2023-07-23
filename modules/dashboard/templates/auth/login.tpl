{% extends "dashboard:layouts/default.tpl" %}

{% block page_title %}Login Screen{% endblock %}


{% block content %}
    <h1>Login Screen</h1>
    <div>
        <form action='/api/v1.0/login?redirect=/dashboard' method='POST'>
            <div class="form-control w-full max-w-xs">
                <label class="label" for='username'>
                    <span class="label-text">Login</span>
                </label>
                <input type="text" name="username" placeholder="john@doe.com" class="input input-bordered w-full max-w-xs" />
            </div>

            <div class="form-control w-full max-w-xs">
                <label class="label" for='password'>
                    <span class="label-text">Password</span>
                </label>
                <input type="password" name="password" class="input input-bordered w-full max-w-xs" />
            </div>

            <input type='submit' value='Login' class="btn"  />
        </form>
    </div>
{% endblock %}