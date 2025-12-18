package problems

// ProblemSeed represents the initial problem library data
type ProblemSeed struct {
	Slug        string
	Title       string
	Description string
	Difficulty  string // "easy", "medium", "hard"
	Topic       string // "arrays", "linked-lists", "trees", etc.
	Tags        []string
}

// SeedData returns the curated initial problem library (21 problems)
func SeedData() []ProblemSeed {
	return []ProblemSeed{
		// Arrays (6 problems)
		{
			Slug:        "two-sum",
			Title:       "Two Sum",
			Description: "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.",
			Difficulty:  "easy",
			Topic:       "arrays",
			Tags:        []string{"hash-table", "two-pointers"},
		},
		{
			Slug:        "best-time-to-buy-sell-stock",
			Title:       "Best Time to Buy and Sell Stock",
			Description: "You are given an array prices where prices[i] is the price of a given stock on the ith day. Maximize profit by buying low and selling high once.",
			Difficulty:  "easy",
			Topic:       "arrays",
			Tags:        []string{"dynamic-programming", "greedy"},
		},
		{
			Slug:        "container-with-most-water",
			Title:       "Container With Most Water",
			Description: "Given n non-negative integers representing vertical lines, find two lines that together with x-axis form container with max water.",
			Difficulty:  "medium",
			Topic:       "arrays",
			Tags:        []string{"two-pointers", "greedy"},
		},
		{
			Slug:        "product-of-array-except-self",
			Title:       "Product of Array Except Self",
			Description: "Given an integer array nums, return array answer such that answer[i] equals product of all elements except nums[i].",
			Difficulty:  "medium",
			Topic:       "arrays",
			Tags:        []string{"prefix-sum", "arrays"},
		},
		{
			Slug:        "maximum-subarray",
			Title:       "Maximum Subarray",
			Description: "Given an integer array nums, find the contiguous subarray with the largest sum and return its sum.",
			Difficulty:  "medium",
			Topic:       "arrays",
			Tags:        []string{"dynamic-programming", "divide-and-conquer"},
		},
		{
			Slug:        "trapping-rain-water",
			Title:       "Trapping Rain Water",
			Description: "Given n non-negative integers representing elevation map, compute how much water can be trapped after raining.",
			Difficulty:  "hard",
			Topic:       "arrays",
			Tags:        []string{"two-pointers", "stack", "dynamic-programming"},
		},

		// Linked Lists (4 problems)
		{
			Slug:        "reverse-linked-list",
			Title:       "Reverse Linked List",
			Description: "Given the head of a singly linked list, reverse the list and return the reversed list.",
			Difficulty:  "easy",
			Topic:       "linked-lists",
			Tags:        []string{"recursion", "iteration"},
		},
		{
			Slug:        "merge-two-sorted-lists",
			Title:       "Merge Two Sorted Lists",
			Description: "Merge two sorted linked lists and return it as a sorted list. The list should be made by splicing together nodes of the first two lists.",
			Difficulty:  "easy",
			Topic:       "linked-lists",
			Tags:        []string{"recursion", "two-pointers"},
		},
		{
			Slug:        "linked-list-cycle",
			Title:       "Linked List Cycle",
			Description: "Given head of a linked list, determine if the linked list has a cycle in it. Use Floyd's Cycle Detection.",
			Difficulty:  "medium",
			Topic:       "linked-lists",
			Tags:        []string{"two-pointers", "floyd-cycle"},
		},
		{
			Slug:        "merge-k-sorted-lists",
			Title:       "Merge K Sorted Lists",
			Description: "You are given an array of k linked-lists, each sorted in ascending order. Merge all into one sorted list.",
			Difficulty:  "hard",
			Topic:       "linked-lists",
			Tags:        []string{"heap", "divide-and-conquer", "priority-queue"},
		},

		// Trees (4 problems)
		{
			Slug:        "invert-binary-tree",
			Title:       "Invert Binary Tree",
			Description: "Given the root of a binary tree, invert the tree and return its root (swap left and right children recursively).",
			Difficulty:  "easy",
			Topic:       "trees",
			Tags:        []string{"recursion", "dfs", "bfs"},
		},
		{
			Slug:        "maximum-depth-of-binary-tree",
			Title:       "Maximum Depth of Binary Tree",
			Description: "Given the root of a binary tree, return its maximum depth (number of nodes along longest path from root to leaf).",
			Difficulty:  "easy",
			Topic:       "trees",
			Tags:        []string{"dfs", "recursion"},
		},
		{
			Slug:        "validate-binary-search-tree",
			Title:       "Validate Binary Search Tree",
			Description: "Given the root of a binary tree, determine if it is a valid binary search tree (BST).",
			Difficulty:  "medium",
			Topic:       "trees",
			Tags:        []string{"dfs", "bst", "recursion"},
		},
		{
			Slug:        "binary-tree-maximum-path-sum",
			Title:       "Binary Tree Maximum Path Sum",
			Description: "Path is sequence of nodes where each pair of adjacent nodes has edge. Path sum is sum of node values. Find maximum.",
			Difficulty:  "hard",
			Topic:       "trees",
			Tags:        []string{"dfs", "recursion", "tree-traversal"},
		},

		// Graphs (3 problems)
		{
			Slug:        "number-of-islands",
			Title:       "Number of Islands",
			Description: "Given m x n 2D grid of '1's (land) and '0's (water), return number of islands. Island is surrounded by water, formed by connecting adjacent lands.",
			Difficulty:  "medium",
			Topic:       "graphs",
			Tags:        []string{"dfs", "bfs", "union-find"},
		},
		{
			Slug:        "clone-graph",
			Title:       "Clone Graph",
			Description: "Given a reference of a node in a connected undirected graph, return a deep copy (clone) of the graph.",
			Difficulty:  "medium",
			Topic:       "graphs",
			Tags:        []string{"dfs", "bfs", "hash-table"},
		},
		{
			Slug:        "course-schedule",
			Title:       "Course Schedule",
			Description: "There are numCourses labeled 0 to n-1. Given prerequisites array, return true if you can finish all courses (detect cycle in directed graph).",
			Difficulty:  "medium",
			Topic:       "graphs",
			Tags:        []string{"topological-sort", "dfs", "bfs"},
		},

		// Sorting (2 problems)
		{
			Slug:        "merge-intervals",
			Title:       "Merge Intervals",
			Description: "Given array of intervals where intervals[i] = [start_i, end_i], merge all overlapping intervals.",
			Difficulty:  "medium",
			Topic:       "sorting",
			Tags:        []string{"sorting", "intervals"},
		},
		{
			Slug:        "sort-colors",
			Title:       "Sort Colors",
			Description: "Given array nums with n objects colored red (0), white (1), blue (2), sort in-place using one-pass Dutch National Flag algorithm.",
			Difficulty:  "medium",
			Topic:       "sorting",
			Tags:        []string{"two-pointers", "dutch-flag", "sorting"},
		},

		// Searching (2 problems)
		{
			Slug:        "binary-search",
			Title:       "Binary Search",
			Description: "Given sorted array nums and target value, return index of target if it exists, otherwise return -1. O(log n) runtime.",
			Difficulty:  "easy",
			Topic:       "searching",
			Tags:        []string{"binary-search", "divide-and-conquer"},
		},
		{
			Slug:        "search-in-rotated-sorted-array",
			Title:       "Search in Rotated Sorted Array",
			Description: "Sorted array nums is possibly rotated at unknown pivot. Given target value, return its index or -1. O(log n) runtime required.",
			Difficulty:  "medium",
			Topic:       "searching",
			Tags:        []string{"binary-search", "arrays"},
		},
	}
}
