{%if true %}<label for="{{ input.Id }}"{% endif %}
    {% if label.Class %}class="{{ label.Class }}" {% endif %}
    {% if style.Style %}class="{{ style.Style }}" {% endif %}
>{{ label.Value }}</label>