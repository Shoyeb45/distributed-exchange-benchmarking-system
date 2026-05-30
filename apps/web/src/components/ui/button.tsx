"use client";

import { Slot } from "@radix-ui/react-slot";
import {
  type ButtonHTMLAttributes,
  forwardRef,
  type ReactElement,
} from "react";

type ButtonVariant = "primary" | "secondary" | "ghost";
type ButtonSize = "sm" | "md" | "lg";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant;
  size?: ButtonSize;
  asChild?: boolean;
};

const variantClass: Record<ButtonVariant, string> = {
  primary: "btn-primary",
  secondary: "btn-secondary",
  ghost: "btn-ghost",
};

const sizeClass: Record<ButtonSize, string> = {
  sm: "btn-sm",
  md: "btn-md",
  lg: "btn-lg",
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className = "",
      variant = "primary",
      size = "md",
      asChild = false,
      type = "button",
      ...props
    },
    ref,
  ): ReactElement => {
    const Comp = asChild ? Slot : "button";
    const classes = ["btn", variantClass[variant], sizeClass[size], className]
      .filter(Boolean)
      .join(" ");

    if (asChild) {
      return <Comp ref={ref} className={classes} {...props} />;
    }

    return <Comp ref={ref} type={type} className={classes} {...props} />;
  },
);

Button.displayName = "Button";
