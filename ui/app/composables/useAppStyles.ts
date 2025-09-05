/**
 * Composable for consistent styling across the application
 * Provides reusable functions for badges, accents, and visual indicators
 */
export const useAppStyles = () => {
  /**
   * Get status badge styling classes
   * @param status - The status string (active, inactive, etc.)
   * @returns Tailwind CSS classes for status badges
   */
  const getStatusBadgeClass = (status: string) => {
    return status === "active"
      ? "bg-chart-2/10 text-chart-2 border border-chart-2/20"
      : "bg-muted text-muted-foreground border border-border"
  }

  /**
   * Get team-based accent color for cards
   * @param team - The team name
   * @returns Accent color key for Card component
   */
  const getTeamAccent = (team: string): "none" | "primary" | "chart-1" | "chart-2" | "chart-3" | "chart-4" => {
    switch (team) {
      case "Customer Success":
        return "chart-2"
      case "Marketing": 
        return "chart-1"
      case "Engineering":
        return "chart-3"
      default:
        return "primary"
    }
  }

  /**
   * Get model badge styling classes
   * @param model - The model name (gpt-4, gpt-3.5-turbo, etc.)
   * @returns Tailwind CSS classes for model badges
   */
  const getModelBadgeClass = (model: string) => {
    return model === "gpt-4"
      ? "bg-primary/10 text-primary border border-primary/20"
      : "bg-chart-4/10 text-chart-4 border border-chart-4/20"
  }

  /**
   * Get activity indicator color based on request count
   * @param requestCount - Number of requests
   * @returns Tailwind CSS class for activity dots
   */
  const getActivityIndicatorClass = (requestCount: number) => {
    return requestCount > 20000 
      ? 'bg-chart-2' 
      : requestCount > 5000 
        ? 'bg-chart-1' 
        : 'bg-muted'
  }

  /**
   * Get priority/urgency color based on value
   * @param value - Numeric value to evaluate
   * @param highThreshold - High priority threshold
   * @param mediumThreshold - Medium priority threshold
   * @returns Chart color key
   */
  const getPriorityColor = (value: number, highThreshold: number, mediumThreshold: number) => {
    return value > highThreshold 
      ? 'chart-2' 
      : value > mediumThreshold 
        ? 'chart-1' 
        : 'muted'
  }

  return {
    getStatusBadgeClass,
    getTeamAccent,
    getModelBadgeClass,
    getActivityIndicatorClass,
    getPriorityColor
  }
}