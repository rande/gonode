declare module 'set-value' {
	function set(target: any, path: string, value: any, options: any): any

	export default set
}

declare module 'get-value' {
	function get(target: any, path: string, options: any): any

	export default get
}
