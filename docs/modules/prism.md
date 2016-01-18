Prism
=====

This plugin allow to render a node with different output depends on the context. The context can be the requested format
or any value.

Each context must be defined and there is one rendered per context. A rendered is a template or a callback.


Configuration
-------------

```toml
[prism]
    [prism.template]


html_desktop => 'reference_to_a_template'
html_mobile => 'reference_to_another_template'
xml_seo => 'reference_to_an_xml_template'

Open questions
--------------

 - How to detect context ?

Usage
-----

  /prism/:uuid => should render the node
  
  